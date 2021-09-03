// Copyright 2021 Deno Land Inc. All rights reserved. MIT License.
package client

import (
	"time"
)

// A User of the Deploy platform.
type User struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl"`
	GitHubID  int       `json:"githubId"`
	IsAdmin   bool      `json:"isAdmin"`
	IsBlocked bool      `json:"isBlocked"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// CurrentUser returns the currently signed in User, the owner of the Token used
// by the client sdk.
func (c *Client) CurrentUser() (User, error) {
	result := User{}
	err := c.request("GET", "/api/user", nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
