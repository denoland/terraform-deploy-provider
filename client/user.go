package client

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl"`
	GitHubID  int       `json:"githubId"`
	IsAdmin   bool      `json:"isAdmin"`
	IsBlocked bool      `json:"isBlocked"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Client) CurrentUser() (User, error) {
	result := User{}
	err := c.request("GET", "/api/user", nil, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
