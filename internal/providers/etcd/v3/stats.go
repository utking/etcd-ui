package v3

import (
	"context"
	"errors"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
)

func (c *Client) ClusterStats() (*types.ClusterStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	statsResp, err := c.client.Cluster.MemberList(ctx)

	if err != nil {
		return nil, err
	}

	return &types.ClusterStats{
		Members:   statsResp.Members,
		MemberID:  statsResp.Header.MemberId,
		ClusterID: statsResp.Header.ClusterId,
		RaftTerm:  statsResp.Header.RaftTerm,
	}, nil
}

func (c *Client) GetLeader() (string, error) {
	var (
		errList []error
		epMap   = make(map[uint64]string, 0)
	)

	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	membersResponse, errMembers := c.client.MemberList(ctx)
	if errMembers != nil {
		return "", nil
	}

	// First, since we can't find the leader with one call, get the members into a map
	for _, member := range membersResponse.Members {
		epMap[member.ID] = member.ClientURLs[0]
	}
	// Now get the endpoints and find which one is the leader
	for _, ep := range c.client.Endpoints() {
		statsResp, err := c.client.Maintenance.Status(ctx, ep)

		if err != nil {
			errList = append(errList, err)

			continue
		}

		if statsResp.Leader == statsResp.Header.MemberId {
			// Now get the URL
			if epURL, exists := epMap[statsResp.Leader]; exists && epURL != "" {
				return epURL, nil
			}
		}
	}

	if len(errList) > 0 {
		return "", errors.Join(errList...)
	}

	return "", nil
}

func (c *Client) EndpointsStatus() (map[uint64]types.EndpointStatusRecord, error) {
	var (
		status  = make(map[uint64]types.EndpointStatusRecord, 0)
		errList []error
	)

	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	for _, ep := range c.client.Endpoints() {
		statsResp, err := c.client.Maintenance.Status(ctx, ep)

		if err != nil {
			errList = append(errList, err)

			continue
		}

		status[statsResp.Header.MemberId] = types.EndpointStatusRecord{
			Errors:   statsResp.Errors,
			DBSize:   statsResp.DbSize,
			IsMaster: statsResp.Leader == statsResp.Header.MemberId,
			Version:  statsResp.Version,
		}
	}

	if len(errList) > 0 {
		return nil, errors.Join(errList...)
	}

	return status, nil
}

func (c *Client) EnableAuth(enable bool) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	if enable {
		_, err = c.client.Auth.AuthEnable(ctx)

		return
	}

	_, err = c.client.Auth.AuthDisable(ctx)

	return
}
