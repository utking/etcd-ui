package v3

import (
	"context"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	"go.etcd.io/etcd/api/v3/authpb"
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

func (c *Client) GetRoles() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	items, err := c.client.RoleList(ctx)

	if err != nil {
		return nil, err
	}

	return items.Roles, nil
}
