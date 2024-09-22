package cluster

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
)

func indexUsers(c echo.Context) error {
	var (
		code  = http.StatusOK
		users = make([]types.UserRecord, 0)
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		users, err = etcdClient.GetUsers()
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"users/list.html",
		map[string]interface{}{
			"Title": "Users",
			"Error": utils.ErrorMessage(err),
			"Items": users,
		},
	)
}

func infoUser(c echo.Context) error {
	var (
		code     = http.StatusOK
		username = c.Param("name")
		userInfo types.UserInfo
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		userInfo, err = etcdClient.UserInfo(username)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"users/info.html",
		map[string]interface{}{
			"Title": "User Info - " + username,
			"Error": utils.ErrorMessage(err),
			"Info":  userInfo,
			"csrf":  c.Get("csrf").(string),
		},
	)
}
