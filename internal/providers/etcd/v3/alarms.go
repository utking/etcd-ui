package v3

import (
	"context"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

func (c *Client) GetAlarms() ([]types.AlarmRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opTimeout)
	defer cancel()

	response, err := c.client.AlarmList(ctx)

	if err != nil {
		return nil, err
	}

	var result = make([]types.AlarmRecord, 0, len(response.Alarms))

	for _, item := range response.Alarms {
		var typeStr string

		switch item.Alarm {
		case etcdserverpb.AlarmType_NOSPACE:
			typeStr = alarmNoSpace
		case etcdserverpb.AlarmType_CORRUPT:
			typeStr = alartCorrupt
		case etcdserverpb.AlarmType_NONE:
			typeStr = alarmNone
		default:
			typeStr = alarmNone
		}

		result = append(result, types.AlarmRecord{
			MemberID: item.MemberID,
			Type:     typeStr,
		})
	}

	return result, nil
}
