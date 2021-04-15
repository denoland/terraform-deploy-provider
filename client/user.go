package client

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	GitHubID  int       `json:"github_id"`
	IsAdmin   bool      `json:"is_admin"`
	IsBlocked bool      `json:"is_blocked"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Client) CurrentUser() (User, error) {
	result := User{}
	err := c.request("GET", "/api/user", nil, nil, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
