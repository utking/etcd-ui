package cluster

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
)

func indexRoles(c echo.Context) error {
	var (
		code  = http.StatusOK
		roles = make([]string, 0)
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		roles, err = etcdClient.GetRoles()
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"roles/list.html",
		map[string]interface{}{
			"Title": "Roles",
			"Error": utils.ErrorMessage(err),
			"Items": roles,
		},
	)
}

func infoRole(c echo.Context) error {
	var (
		code     = http.StatusOK
		role     = c.Param("name")
		roleInfo types.RoleInfo
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		roleInfo, err = etcdClient.RoleInfo(role)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"roles/info.html",
		map[string]interface{}{
			"Title": "Role Info - " + role,
			"Error": utils.ErrorMessage(err),
			"Info":  roleInfo,
			"csrf":  c.Get("csrf").(string),
		},
	)
}
