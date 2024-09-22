package cluster

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
	"github.com/utking/etcd-ui/internal/types/requests"
)

func indexRoles(c echo.Context) error {
	var (
		code   = http.StatusOK
		roles  = make([]string, 0)
		filter = c.QueryParam("filter")
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		roles, err = etcdClient.GetRoles(filter)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"roles/list.html",
		map[string]interface{}{
			"Title":  "Roles",
			"Error":  utils.ErrorMessage(err),
			"Items":  roles,
			"Filter": filter,
		},
	)
}

func infoRole(c echo.Context) error {
	var (
		code     = http.StatusOK
		role     = c.QueryParam("name")
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

func createRole(c echo.Context) error {
	var (
		code       = http.StatusOK
		item       requests.RoleCreate
		err        error
		etcdClient *v3.Client
	)

	etcdClient, err = v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		// GET request
		item.Name = c.Param("name")
		if c.Request().Method == http.MethodGet && item.Name != "" {
			roleItem, getErr := etcdClient.RoleInfo(item.Name)
			if getErr == nil {
				readPerms := make([]requests.KVPerm, 0, len(roleItem.KVRead))
				writePerms := make([]requests.KVPerm, 0, len(roleItem.KVWrite))

				for _, permItem := range roleItem.KVRead {
					readPerms = append(readPerms, requests.KVPerm{
						Key:     permItem.Key,
						IsRange: permItem.RangeEnd != "",
					})
				}

				for _, permItem := range roleItem.KVWrite {
					writePerms = append(writePerms, requests.KVPerm{
						Key:     permItem.Key,
						IsRange: permItem.RangeEnd != "",
					})
				}

				item.Name = roleItem.Name
				item.ReadPerms = readPerms
				item.WritePerms = writePerms
			}
		} else if c.Request().Method == http.MethodPost {
			// POST request
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddRole(item.Name, item.GetReadPerms(), item.GetWritePerms())

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/roles/list?filter="+item.Name)
				}
			}
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"roles/create.html",
		map[string]interface{}{
			"Title": "Create Role", // FIX: for edit
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func deleteRole(c echo.Context) error {
	var (
		name       = c.FormValue("name")
		method     = c.FormValue("_method")
		err        error
		etcdClient *v3.Client
	)

	if c.Request().Method == http.MethodPost && strings.EqualFold(method, http.MethodDelete) {
		etcdClient, err = v3.New(
			utils.GetEntryPoints(),
			utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
			utils.GetUsername(), utils.GetPassword(),
		)

		if err == nil {
			_ = etcdClient.DeleteRole(name)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/roles/list?filter="+name)
}
