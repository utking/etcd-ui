package v3

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	"go.etcd.io/etcd/api/v3/authpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (c *Client) RoleInfo(roleName string) (types.RoleInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	response, err := c.client.RoleGet(ctx, roleName)

	if err != nil {
		return types.RoleInfo{}, err
	}

	var (
		kvPerms = make([]types.KVPerm, 0, len(response.Perm))
	)

	for _, perm := range response.Perm {
		kvPerm := types.KVPerm{Key: string(perm.Key), RangeEnd: string(perm.RangeEnd)}

		switch perm.PermType {
		case authpb.READWRITE:
			kvPerm.Type = types.PermReadWrite
		case authpb.READ:
			kvPerm.Type = types.PermRead
		case authpb.WRITE:
			kvPerm.Type = types.PermWrite
		default:
			continue
		}

		kvPerms = append(kvPerms, kvPerm)
	}

	return types.RoleInfo{
		Name:  roleName,
		Perms: kvPerms,
	}, nil
}

func (c *Client) GetRoles(filter string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	items, err := c.client.RoleList(ctx)

	if err != nil {
		return nil, err
	}

	if filter == "" {
		return items.Roles, nil
	}

	var filtered []string

	for _, role := range items.Roles {
		if strings.Contains(role, filter) {
			filtered = append(filtered, role)
		}
	}

	return filtered, nil
}

// AddRole creates a role (if does not exist) and assigns it the required permissions
func (c *Client) AddRole(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	existingRole, _ := c.client.RoleGet(ctx, name)

	if existingRole != nil {
		return fmt.Errorf("the role named %q already exists", name)
	}

	created, err := c.client.RoleAdd(ctx, name)
	if err != nil {
		return err
	}

	if created == nil {
		return fmt.Errorf("the role could not be create due to an unknown error")
	}

	return nil
}

func (c *Client) DeleteRole(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.RoleDelete(ctx, name)

	return err
}

func (c *Client) RevokePermissions(name string, perms []types.KVPerm) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	var errList []error

	for _, perm := range perms {
		_, revokeErr := c.client.RoleRevokePermission(ctx, name, perm.Key, perm.RangeEnd)
		if revokeErr != nil {
			errList = append(errList, revokeErr)
		}
	}

	return errors.Join(errList...)
}

func (c *Client) GrantPermissions(name string, perms []types.KVPerm) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	var errList []error

	for _, perm := range perms {
		rangeEnd := perm.RangeEnd
		// Define the range-end if we grant with prefix
		if perm.RangeEnd == "" && perm.IsRange {
			rangeEnd = clientv3.GetPrefixRangeEnd(perm.Key)
		}

		_, grantErr := c.client.RoleGrantPermission(
			ctx, name,
			perm.Key, rangeEnd, clientv3.PermissionType(perm.Type),
		)
		if grantErr != nil {
			errList = append(errList, grantErr)
		}
	}

	return errors.Join(errList...)
}
