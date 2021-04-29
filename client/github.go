// Copyright 2021 William Perron. All rights reserved. MIT License.
package client

import (
	"bytes"
	"encoding/json"
)

// LinkProjectRequest is the expected request body schema for the LinkProject
// function.
type LinkProjectRequest struct {
	ProjectID    string `json:"projectId"`
	Organization string `json:"organization"`
	Repo         string `json:"repo"`
	Entrypoint   string `json:"entrypoint"`
}

func (c *Client) LinkProject(req LinkProjectRequest) (Project, error) {
	bs, err := json.Marshal(req)
	if err != nil {
		return Project{}, nil
	}

	res := Project{}
	err = c.request("POST", "/api/github/link", nil, bytes.NewBuffer(bs), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
