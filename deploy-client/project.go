package deployclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Project struct {
	Id                      string
	Name                    string      `json:"name"`
	Git                     GitHubLink  `json:"git,omitempty"`
	ProductionDeployment    *Deployment `json:"productionDeployment,omitempty"`
	HasProductionDeployment bool        `json:"hasProductionDeployment"`
	EnvVars                 EnvVars     `json:"envVars"`
	UpdatedAt               time.Time   `json:"updatedAt"`
	CreatedAt               time.Time   `json:"createdAt"`
}

type EnvVars map[string]string

type GitHubLink struct {
	Repository Repository `json:"repository"`
	Entrypoint string     `json:"entrypoint"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type Repository struct {
	Id    int    `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type Deployment struct {
	Id             string          `json:"id"`
	Url            string          `json:"url"`
	DomainMappings []DomainMapping `json:"domainMappings"`
	RelatedCommit  CommitInfo      `json:"relatedCommit"`
	Project        *Project        `json:"project"`
	ProjectId      string          `json:"projectId"`
	EnvVars        EnvVars         `json:"envVars"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type DomainMapping struct {
	Domain    string    `json:"domain"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommitInfo struct {
	Hash                 string `json:"hash"`
	Message              string `json:"message"`
	AuthorName           string `json:"authorName"`
	AuthorEmail          string `json:"authorEmail"`
	AuthorGitHubUsername string `json:"authorGithubUsername,omitempty"`
	Url                  string `json:"url,omitempty"`
}

type Domain struct {
	Domain       string    `json:"domain"`
	Token        string    `json:"token"`
	IsValidated  bool      `json:"isValidated"`
	Certificates []string  `json:"certificates"` // TODO(wperron) implement TlsCipher struct
	ProjectId    string    `json:"projectId"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

const (
	TlsCipherRsa = "rsa"
	TlsCipherEc  = "ec"
)

func (c *Client) ListProjects() ([]Project, error) {
	result := []Project{}
	err := c.request("GET", "/api/projects", nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

type CreateProjectRequest struct {
	Name    string  `json:"name"`
	EnvVars EnvVars `json:"envVars"`
}

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

type UpdateProjectRequest struct {
	Name string `json:"name"`
}

func (c *Client) UpdateProject(projectId string, newName string) error {
	path := fmt.Sprintf("/api/projects/%s", projectId)
	project := UpdateProjectRequest{
		Name: newName,
	}

	bs, err := json.Marshal(project)
	if err != nil {
		return err
	}

	return c.request("PATCH", path, nil, bytes.NewBuffer(bs), nil)
}

func (c *Client) DeleteProject(projectId string) error {
	path := fmt.Sprintf("/api/projects/%s", projectId)
	return c.request("DELETE", path, nil, nil, nil)
}

func (c *Client) GetProject(projectId string) (Project, error) {
	path := fmt.Sprintf("/api/projects/%s", projectId)
	result := Project{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

type NewDeploymentRequest struct {
	Url        string `json:"url"`
	Production bool   `json:"production,omitempty"`
}

func (c *Client) NewProjectDeployment(projectId string, depl NewDeploymentRequest) (Deployment, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments", projectId)

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

func (c *Client) ListDeployments(projectId string, pageOpts PageOptions) ([]Deployment, PagingInfo, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments", projectId)
	// expected []Deployment at position 0 and PagingInfo at position 1
	result := []interface{}{}

	qs := url.Values{}
	if pageOpts.Page != 0 {
		qs.Set("page", fmt.Sprint(pageOpts.Page))
	}
	if pageOpts.Limit != 0 {
		qs.Set("limit", fmt.Sprint(pageOpts.Limit))
	}

	err := c.request("GET", path, qs, nil, result)
	if err != nil {
		return []Deployment{}, PagingInfo{}, err
	}

	return result[0].([]Deployment), result[1].(PagingInfo), nil
}

func (c *Client) GetDeployment(projectId string, deploymentId string) (Deployment, error) {
	path := fmt.Sprintf("/api/projects/%s/deployments/%s", projectId, deploymentId)
	result := Deployment{}
	err := c.request("GET", path, nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) GetLogs(projectId string, deploymentId string) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (c *Client) UpdateEnvVars(projectId string, newVars EnvVars) error {
	path := fmt.Sprintf("/api/projects/%s/env", projectId)

	bs, err := json.Marshal(newVars)
	if err != nil {
		return err
	}

	return c.request("POST", path, nil, bytes.NewBuffer(bs), nil)
}

func (c *Client) Unlink(projectId string) error {
	path := fmt.Sprintf("/api/projects/%s/git", projectId)
	return c.request("DELETE", path, nil, nil, nil)
}

func (c *Client) ListDomains(projectId string) ([]Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains", projectId)
	result := []Domain{}
	err := c.request("GET", path, nil, nil, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

type AddDomainRequest struct {
	Domain Domain `json:"domain"`
}

func (c *Client) AddDomain(projectId string, newDomain AddDomainRequest) (Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains", projectId)

	bs, err := json.Marshal(newDomain)
	if err != nil {
		return Domain{}, err
	}

	result := Domain{}
	err = c.request("POST", path, nil, bytes.NewBuffer(bs), result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) GetDomain(projectId, domainName string) (Domain, error) {
	path := fmt.Sprintf("/api/projects/%s/domains/%s", projectId, domainName)
	result := Domain{}
	err := c.request("GET", path, nil, nil, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) DeleteDomain(projectId, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s", projectId, domainName)
	return c.request("DELETE", path, nil, nil, nil)
}

func (c *Client) VerifyDomain(projectId, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s/verify", projectId, domainName)
	return c.request("POST", path, nil, nil, nil)
}

func (c *Client) ProvisionCertificate(projectId, domainName string) error {
	path := fmt.Sprintf("/api/projects/%s/domains/%s/certificates", projectId, domainName)
	return c.request("POST", path, nil, nil, nil)
}
