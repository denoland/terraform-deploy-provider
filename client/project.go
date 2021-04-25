// Copyright 2021 William Perron. All rights reserved. MIT License.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

// A Project represents a Deploy project.
//
// Projects can be linked to a public GitHub repository or point directly to a
// publicly accessible URL. The Project resource contains certain global
// settings like the name of the project, domain names and environment variables.
// Projects are made up of Deployments. Every time the source URL is updated or
// a commit is pushed to the default branch on GitHub, a new Deployment is
// created and the 'production' deployment is updated.
//
// The ProductionDeployment property is a pointer to the Deployment object that
// represents the latest version of the project. When a project is first created
// this value is nil and HasProductionDeployment is set to `false`.
//
// The Git property is only set if the project is linked to a GitHub repository,
// otherwise it is nil.
//
// A Project also has a circular reference to a Deployment. The Deployment
// property is only set when accessing the Project directly in the API,
// otherwise it is omitted.
type Project struct {
	ID                      string      `json:"id"`
	Name                    string      `json:"name"`
	Git                     *GitHubLink `json:"git,omitempty"`
	ProductionDeployment    *Deployment `json:"productionDeployment,omitempty"`
	HasProductionDeployment bool        `json:"hasProductionDeployment"`
	EnvVars                 EnvVars     `json:"envVars"`
	UpdatedAt               string      `json:"updatedAt"`
	CreatedAt               string      `json:"createdAt"`
}

// EnvVars is a simple map used to created environment variables in a Project.
type EnvVars map[string]string

// A GitHubLink is used in a Project to link it to a GitHub repository. it
// contains a Repository struct and an entrypoint corresponding to the source
// code file used as the entrypoint of the project.
type GitHubLink struct {
	Repository Repository `json:"repository"`
	Entrypoint string     `json:"entrypoint"`
	UpdatedAt  string     `json:"updatedAt"`
	CreatedAt  string     `json:"createdAt"`
}

// A Repository is a simple structure containing the information identifying the
// GitHub repository.
type Repository struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// A Deployment is an immutable version of a Project's source code.
//
// Each Deployment contains a URL to the source code's entrypoint, the auto
// generated domain name that was created for the Deployment and a copy of the
// environment variables of the project when the Deployment was created.
//
// If the project is linked to a GitHub repository, it will also contain a
// CommitInfo containing the summary of the commit.
//
// A Deployment also has a circular reference to a Project. The Project property
// is only set when accessing the Deployment directly in the API, otherwise it
// is omitted.
type Deployment struct {
	ID             string          `json:"id"`
	URL            string          `json:"url"`
	DomainMappings []DomainMapping `json:"domainMappings"`
	RelatedCommit  *CommitInfo     `json:"relatedCommit,omitempty"`
	Project        *Project        `json:"project"`
	ProjectID      string          `json:"projectId"`
	EnvVars        EnvVars         `json:"envVars"`
	UpdatedAt      string          `json:"updatedAt"`
	CreatedAt      string          `json:"createdAt"`
}

// A DomainMapping is a simple struct containing to immutable domain name of a
// Deployment.
type DomainMapping struct {
	Domain    string `json:"domain"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// A CommitInfo is used for Projects linked to a GitHub repository. It contains
// the information about the commit that triggered a new Deployment for the
// Project.
type CommitInfo struct {
	Hash                 string `json:"hash"`
	Message              string `json:"message"`
	AuthorName           string `json:"authorName"`
	AuthorEmail          string `json:"authorEmail"`
	AuthorGitHubUsername string `json:"authorGithubUsername,omitempty"`
	URL                  string `json:"url,omitempty"`
}

// A Domain is a custom domain name for a Project.
type Domain struct {
	Domain       string   `json:"domain"`
	Token        string   `json:"token"`
	IsValidated  bool     `json:"isValidated"`
	Certificates []string `json:"certificates"` // TODO(wperron) implement TlsCipher struct
	ProjectID    string   `json:"projectId"`
	UpdatedAt    string   `json:"updatedAt"`
	CreatedAt    string   `json:"createdAt"`
}

// Possible values for the Certificates property of the Domain struct
const (
	TLSCipherRsa = "rsa"
	TLSCipherEc  = "ec"
)

// ListProjects returns a slice of Projects owned by the current User.
func (c *Client) ListProjects() ([]Project, error) {
	result := []Project{}
	err := c.request("GET", "/api/projects", nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// CreateProjectRequest is the expected request body schema for the
// CreateProject function.
type CreateProjectRequest struct {
	Name    string  `json:"name"`
	EnvVars EnvVars `json:"envVars"`
}

// CreateProject creates a new project with the given name.
func (c *Client) CreateProject(name string, envVars EnvVars) (Project, error) {
	project := CreateProjectRequest{
		Name:    name,
		EnvVars: envVars,
	}

	bs, err := json.Marshal(project)
	if err != nil {
		return Project{}, err
	}

	res := Project{}
	err = c.request("POST", "/api/projects", nil, bytes.NewBuffer(bs), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateProjectRequest is the expected request body schema for the
// UpdateProject function.
type UpdateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateProject modifies an existing Project.
func (c *Client) UpdateProject(projectID string, newName string) error {
	path := fmt.Sprintf("/api/projects/%s", projectID)
	project := UpdateProjectRequest{
		Name: newName,
	}

	bs, err := json.Marshal(project)
	if err != nil {
		return err
	}

	return c.request("PATCH", path, nil, bytes.NewBuffer(bs), nil)
}

// DeleteProject deletes a Project and all of its associated Deployments.
func (c *Client) DeleteProject(projectID string) error {
	path := fmt.Sprintf("/api/projects/%s", projectID)
	return c.request("DELETE", path, nil, nil, nil)
}

// GetProject returns the information about a given Project.
func (c *Client) GetProject(projectID string) (Project, error) {
	path := fmt.Sprintf("/api/projects/%s", projectID)
	result := Project{}
	log.Printf("[DEBUG] GET %s", path)
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		log.Printf("[DEBUG] Could not find project: %s", err)
		return result, err
	}

	bs, _ := json.Marshal(result)
	log.Printf("[DEBUG] GET Project response body: %s", string(bs))

	return result, nil
}

// NewDeploymentRequest is the expected request body schema for the
// NewProjectDeployment function.
type NewDeploymentRequest struct {
	URL        string `json:"url"`
	Production bool   `json:"production,omitempty"`
}

// NewProjectDeployment creates a new Deployment for a Project. The URL param
// is the URL of the source code used for the deployment. The URL needs to be
// publicly available.
//
// Note that this function is not needed to create a new Deployment for a
// Project linked to a GitHub repository.
func (c *Client) NewProjectDeployment(projectID string, depl NewDeploymentRequest) (Deployment, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments", projectID)

	bs, err := json.Marshal(depl)
	if err != nil {
		return Deployment{}, err
	}

	res := Deployment{}
	err = c.request("POST", path, nil, bytes.NewBuffer(bs), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// ListDeployments returns a slice of Deployments owned by the current User for
// a given Project.
func (c *Client) ListDeployments(projectID string, pageOpts PageOptions) ([]Deployment, PagingInfo, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments", projectID)
	// expected []Deployment at position 0 and PagingInfo at position 1
	result := []interface{}{}

	qs := url.Values{}
	if pageOpts.Page != 0 {
		qs.Set("page", fmt.Sprint(pageOpts.Page))
	}
	if pageOpts.Limit != 0 {
		qs.Set("limit", fmt.Sprint(pageOpts.Limit))
	}

	err := c.request("GET", path, qs, nil, &result)
	if err != nil {
		return []Deployment{}, PagingInfo{}, err
	}

	return result[0].([]Deployment), result[1].(PagingInfo), nil
}

// GetDeployment returns the information about a given Deployment.
func (c *Client) GetDeployment(projectID string, deploymentID string) (Deployment, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments/%s", projectID, deploymentID)
	result := Deployment{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetLogs returns the log lines from a given Deployment
func (c *Client) GetLogs(projectID string, deploymentID string) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

// UpdateEnvVars overwrites the environment variables of a given Project.
func (c *Client) UpdateEnvVars(projectID string, newVars EnvVars) error {
	path := fmt.Sprintf("/api/projects/%s/env", projectID)

	bs, err := json.Marshal(newVars)
	if err != nil {
		return err
	}

	return c.request("POST", path, nil, bytes.NewBuffer(bs), nil)
}

// Unlink removes the GitHub integration of a given project.
//
// This only affects future Deployments. Any active Deployment that was created
// when the GitHub repository was linked will still be active.
func (c *Client) Unlink(projectID string) error {
	path := fmt.Sprintf("/api/projects/%s/git", projectID)
	return c.request("DELETE", path, nil, nil, nil)
}

// ListDomains returns a list of all the custom domain names associated to the
// Project.
func (c *Client) ListDomains(projectID string) ([]Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains", projectID)
	result := []Domain{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// AddDomainRequest is the expected request body schema for the AddDomain
// function.
type AddDomainRequest struct {
	Domain Domain `json:"domain"`
}

// AddDomain adds a custom domain name to the project. This is typically
// followed by the VerifyDomain function
func (c *Client) AddDomain(projectID string, newDomain AddDomainRequest) (Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains", projectID)

	bs, err := json.Marshal(newDomain)
	if err != nil {
		return Domain{}, err
	}

	result := Domain{}
	err = c.request("POST", path, nil, bytes.NewBuffer(bs), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetDomain returns the information about a given custom domain name.
//
// This is typically used to retrieve the information about the different
// records that must be created by the user to properly verify the domain.
func (c *Client) GetDomain(projectID, domainName string) (Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains/%s", projectID, domainName)
	result := Domain{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// DeleteDomain removes a custom domain name associated with a project.
//
// This action only removes the custom domain name resolution on the Deploy side.
// The DNS records will still have to be removed on the user's registrar.
func (c *Client) DeleteDomain(projectID, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s", projectID, domainName)
	return c.request("DELETE", path, nil, nil, nil)
}

// VerifyDomain sends the signal to Deploy to verify the DNS records for custom
// domain names. The DNS records must exist prior to starting the verification
// process. Deploy will not create these for you.
func (c *Client) VerifyDomain(projectID, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s/verify", projectID, domainName)
	return c.request("POST", path, nil, nil, nil)
}

// ProvisionCertificate creates a valid TLS certificate for a custom domain name.
func (c *Client) ProvisionCertificate(projectID, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s/certificates", projectID, domainName)
	return c.request("POST", path, nil, nil, nil)
}
