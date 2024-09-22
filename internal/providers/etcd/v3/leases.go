package v3

import (
	"context"
	"fmt"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Returns grantedTTL, remainingTTL, error
func (c *Client) getLeaseTTL(id clientv3.LeaseID) (granted, current int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	ttlResponse, err := c.client.Lease.TimeToLive(ctx, id)

	if err != nil {
		return 0, 0, err
	}

	return ttlResponse.GrantedTTL, ttlResponse.TTL, nil
}

func (c *Client) GetLeases(filter string) ([]types.LeaseRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	leases, err := c.client.Leases(ctx)

	if err != nil {
		return nil, err
	}

	var result = make([]types.LeaseRecord, 0, 1)

	if filter == "" {
		result = make([]types.LeaseRecord, 0, len(leases.Leases))
	}

	for _, lease := range leases.Leases {
		if filter != "" && filter != fmt.Sprintf("%x", lease.ID) {
			continue
		}

		record := types.LeaseRecord{
			ID: int64(lease.ID), TTL: 0,
			HexID: fmt.Sprintf("%x", lease.ID),
		}
		record.GrantedTTL, record.TTL, _ = c.getLeaseTTL(lease.ID)
		result = append(result, record)
	}

	return result, nil
}

func (c *Client) GrantLease(ttl int64) (types.LeaseRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	leaseResponse, err := c.client.Grant(ctx, ttl)

	if err != nil || leaseResponse == nil {
		return types.LeaseRecord{}, err
	}

	return types.LeaseRecord{
		ID: int64(leaseResponse.ID),
	}, nil
}

func (c *Client) RenewLease(id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	leaseResponse, err := c.client.Lease.KeepAliveOnce(ctx, clientv3.LeaseID(id))

	if err != nil || leaseResponse == nil {
		return false, err
	}

	return true, nil
}

func (c *Client) LeaseInfo(id int64) (types.LeaseRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	leaseResponse, err := c.client.Lease.TimeToLive(ctx, clientv3.LeaseID(id), clientv3.WithAttachedKeys())

	if err != nil || leaseResponse == nil {
		return types.LeaseRecord{}, err
	}

	var assignedKeys = make([]string, 0, len(leaseResponse.Keys))

	for _, key := range leaseResponse.Keys {
		assignedKeys = append(assignedKeys, string(key))
	}

	var result = types.LeaseRecord{
		ID:         id,
		HexID:      fmt.Sprintf("%x", id),
		GrantedTTL: leaseResponse.GrantedTTL,
		TTL:        leaseResponse.TTL,
		Keys:       assignedKeys,
	}

	return result, nil
}

func (c *Client) DeleteLease(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.Lease.Revoke(ctx, clientv3.LeaseID(id))

	return err
}
