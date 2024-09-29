package cluster

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
)

func clusterIndex(c echo.Context) error {
	var (
		code    = http.StatusOK
		stats   *types.ClusterStats
		members = make([]types.MemberRecord, 0)
		alarms  = make([]types.AlarmRecord, 0)
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		stats, err = etcdClient.ClusterStats()
		alarms, _ = etcdClient.GetAlarms()
	}

	if err != nil {
		code = http.StatusInternalServerError
	} else {
		endpointsStatus, statusErr := etcdClient.EndpointsStatus()
		if statusErr == nil {
			for _, member := range stats.Members {
				memberItem := types.MemberRecord{
					ID:         member.ID,
					Name:       member.Name,
					PeerURLs:   member.PeerURLs,
					ClientURLs: member.ClientURLs,
				}

				if epStatus, exists := endpointsStatus[member.ID]; exists {
					memberItem.Health = epStatus
					memberItem.Version = epStatus.Version
				}

				members = append(members, memberItem)
			}
		}
	}

	return c.Render(
		code,
		"cluster/stats.html",
		map[string]interface{}{
			"Title":      "Cluster Stats",
			"Error":      utils.ErrorMessage(err),
			"Items":      members,
			"SingleNode": len(members) == 1,
			"Alarms":     alarms,
			"Header":     stats,
		},
	)
}

func electNewLeader(c echo.Context) error {
	var (
		memberIDParam = c.Param("id")
		electionErr   error
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		memberID, parseErr := strconv.ParseUint(memberIDParam, 10, 64)
		if parseErr == nil {
			_, electionErr = etcdClient.MoveLeader(memberID)
			if electionErr != nil {
				err = electionErr
			}
		}
	}

	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"Error": utils.ErrorMessage(err),
			},
		)
	}

	time.Sleep(time.Second)

	return c.Redirect(http.StatusSeeOther, "/cluster/stats")
}
