package cluster

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
	"github.com/utking/etcd-ui/internal/types/requests"
)

func indexLeases(c echo.Context) error {
	var (
		code   = http.StatusOK
		leases = make([]types.LeaseRecord, 0)
		filter = c.QueryParam("filter")
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		leases, err = etcdClient.GetLeases(filter)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"leases/list.html",
		map[string]interface{}{
			"Title":  "Active Leases",
			"Error":  utils.ErrorMessage(err),
			"Items":  leases,
			"Filter": filter,
		},
	)
}

func infoLease(c echo.Context) error {
	var (
		code      = http.StatusOK
		id        = c.Param("id")
		leaseInfo types.LeaseRecord
	)

	leaseID, _ := strconv.ParseInt(id, 10, 64)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		leaseInfo, err = etcdClient.LeaseInfo(leaseID)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"leases/info.html",
		map[string]interface{}{
			"Title": "Lease Info - " + leaseInfo.HexID,
			"Error": utils.ErrorMessage(err),
			"Info":  leaseInfo,
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func deleteLease(c echo.Context) error {
	var (
		id         = c.FormValue("id")
		method     = c.FormValue("_method")
		err        error
		etcdClient *v3.Client
	)

	leaseID, _ := strconv.ParseInt(id, 10, 64)

	if c.Request().Method == http.MethodPost && strings.EqualFold(method, http.MethodDelete) {
		etcdClient, err = v3.New(
			utils.GetEntryPoints(),
			utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
			utils.GetUsername(), utils.GetPassword(),
		)

		if err == nil {
			_ = etcdClient.DeleteLease(leaseID)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/leases/list")
}

func createLease(c echo.Context) error {
	var (
		code       = http.StatusOK
		id         = c.Param("id")
		item       requests.LeaseCreate
		err        error
		etcdClient *v3.Client
		renewed    bool
		title      = "Create"
	)

	etcdClient, err = v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		item.LeaseID, _ = strconv.ParseInt(id, 10, 64)

		if c.Request().Method == http.MethodGet && item.LeaseID != 0 {
			title = "Renew"

			leaseItem, getErr := etcdClient.LeaseInfo(item.LeaseID)
			if getErr == nil {
				item.TTL = leaseItem.GrantedTTL
			}
		} else if c.Request().Method == http.MethodPost {
			err = c.Bind(&item)
			if err == nil {
				if item.LeaseID > 0 {
					// Renew existing
					title = "Renew"
					renewed, err = etcdClient.RenewLease(item.LeaseID)

					if err == nil && renewed {
						return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/cluster/leases/list?filter=%x", item.LeaseID))
					}
				} else {
					var newLease types.LeaseRecord

					// Create new
					newLease, err = etcdClient.GrantLease(item.TTL)
					if err == nil {
						return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/cluster/leases/list?filter=%x", newLease.ID))
					}
				}
			}
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"leases/create.html",
		map[string]interface{}{
			"Title": title + " Lease",
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}
