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
		kvRead  = make([]types.KVPerm, 0, len(response.Perm))
		kvWrite = make([]types.KVPerm, 0, len(response.Perm))
	)

	for _, perm := range response.Perm {
		kvPerm := types.KVPerm{Key: string(perm.Key), RangeEnd: string(perm.RangeEnd)}

		switch perm.PermType {
		case authpb.READWRITE:
			kvRead = append(kvRead, kvPerm)
			kvWrite = append(kvWrite, kvPerm)
		case authpb.READ:
			kvRead = append(kvRead, kvPerm)
		case authpb.WRITE:
			kvWrite = append(kvWrite, kvPerm)
		default:
			continue
		}
	}

	return types.RoleInfo{
		Name:    roleName,
		KVRead:  kvRead,
		KVWrite: kvWrite,
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
func (c *Client) AddRole(
	name string,
	kvRead, kvWrite []types.KVPerm,
) error {
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

	var errList []error

	// Role created. Grant permissions
	for _, readPerm := range kvRead {
		fmt.Println("Grant read perm for", name)
		_, permErr := c.client.RoleGrantPermission(
			ctx, name,
			readPerm.Key, readPerm.RangeEnd, clientv3.PermissionType(authpb.READ),
		)

		if permErr != nil {
			errList = append(errList, permErr)
		}
	}

	for _, write := range kvWrite {
		fmt.Println("Grant write perm for", name)
		_, permErr := c.client.RoleGrantPermission(
			ctx, name,
			write.Key, write.RangeEnd, clientv3.PermissionType(authpb.WRITE),
		)

		if permErr != nil {
			errList = append(errList, permErr)
		}
	}

	fmt.Println("Done for", name)
	return errors.Join(errList...)
}

func (c *Client) DeleteRole(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.RoleDelete(ctx, name)

	return err
}
