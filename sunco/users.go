package sunco

import (
	"context"
	"encoding/json"
	"fmt"
)

type UsersAPI interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUser(ctx context.Context, userId string) (User, error)
}

func (c *Client) CreateUser(ctx context.Context, user User) (User, error) {
	var response UserResponse

	body, err := c.Post(ctx, "/users", user)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return User{}, err
	}

	return response.User, nil
}

func (c *Client) GetUser(ctx context.Context, userId string) (User, error) {
	var response UserResponse

	body, err := c.Get(ctx, fmt.Sprintf("/users/%s", userId))
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return User{}, err
	}

	return response.User, nil
}
