package v3

import (
	"context"
	"strings"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
)

func (c *Client) UserInfo(username string) (types.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	userResponse, err := c.client.UserGet(ctx, username)

	if err != nil {
		return types.UserInfo{}, err
	}

	return types.UserInfo{
		Name:  types.UserRecord(username),
		Roles: userResponse.Roles,
	}, nil
}

func (c *Client) GetUsers(filter string) ([]types.UserRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	users, err := c.client.UserList(ctx)

	if err != nil {
		return nil, err
	}

	var result = make([]types.UserRecord, 0, len(users.Users))

	for _, user := range users.Users {
		if filter != "" && !strings.Contains(user, filter) {
			continue
		}

		result = append(result, types.UserRecord(user))
	}

	return result, nil
}

func (c *Client) AddUser(username string, password types.SensitiveStr) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.UserAdd(ctx, username, password.Unwrap())

	return err
}

func (c *Client) DeleteUser(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.UserDelete(ctx, name)

	return err
}
