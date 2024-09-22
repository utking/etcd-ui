package v3

import (
	"context"
	"fmt"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (c *Client) Get(key string) (*types.KVItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	resp, err := c.client.Get(ctx, key)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("the key %q has not been found", key)
	}

	_, ttl, _ := c.getLeaseTTL(clientv3.LeaseID(resp.Kvs[0].Lease))

	return &types.KVItem{
		Key:      string(resp.Kvs[0].Key),
		Value:    types.SensitiveStr(resp.Kvs[0].Value),
		LeaseTTL: ttl,
		LeaseID:  resp.Kvs[0].Lease,
	}, nil
}

func (c *Client) Put(key, val string, ttl int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	var opts = make([]clientv3.OpOption, 0)

	if ttl > 0 {
		resp, _ := c.client.Grant(ctx, ttl)
		opts = append(opts, clientv3.WithLease(resp.ID))
	}

	_, err := c.client.Put(
		ctx,
		key,
		val,
		opts...,
	)

	return err
}

func (c *Client) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	_, err := c.client.Delete(ctx, key)

	return err
}

func (c *Client) Watch(key string) <-chan *clientv3.Event {
	events := make(chan *clientv3.Event)

	wCh := c.client.Watch(
		context.Background(),
		key,
		clientv3.WithPrefix(),
		// clientv3.WithFilterDelete(),
		clientv3.WithPrevKV(),
	)

	go func() {
		defer close(events)

		for wResp := range wCh {
			for _, evt := range wResp.Events {
				events <- evt
			}
		}
	}()

	return events
}

func (c *Client) ListKeys(prefix string, limit int64, withData bool) ([]types.KVItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	var opts = []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithLimit(limit),
	}

	if !withData {
		opts = append(opts, clientv3.WithKeysOnly())
	}

	response, listErr := c.client.KV.Get(ctx, prefix, opts...)

	if listErr != nil {
		return []types.KVItem{}, listErr
	}

	var result = make([]types.KVItem, 0, len(response.Kvs))

	for _, kvItem := range response.Kvs {
		_, ttl, _ := c.getLeaseTTL(clientv3.LeaseID(kvItem.Lease))
		result = append(result, types.KVItem{
			Key:      string(kvItem.Key),
			Value:    types.SensitiveStr(kvItem.Value),
			LeaseTTL: ttl,
		})
	}

	return result, nil
}
