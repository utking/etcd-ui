package v3

import (
	"context"
	"fmt"
)

func (c *Client) MoveLeader(memberID uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	// First, get the leader's endpoint since only the leader can make the move
	leaderEndpoint, leaderErr := c.GetLeader()
	if leaderErr == nil && leaderEndpoint != "" {
		// Set the leader's endpoint
		c.client.SetEndpoints(leaderEndpoint)
		// make the move
		_, err := c.client.MoveLeader(ctx, memberID)

		return err == nil, err
	}

	return false, fmt.Errorf("could not define the current leader")
}
