package cluster

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
		if c.Request().Method == http.MethodPost {
			// POST request
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddRole(item.Name)

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
			"Title": "Create Role",
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

func editRolePermissions(c echo.Context) error {
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
				kvPerms := make([]requests.KVPerm, 0, len(roleItem.Perms))

				for _, permItem := range roleItem.Perms {
					kvPerms = append(kvPerms, requests.KVPerm{
						Key:      permItem.Key,
						RangeEnd: permItem.RangeEnd,
						Type:     permItem.Type,
					})
				}

				item.Name = roleItem.Name
				item.Perms = kvPerms
			}
		} else if c.Request().Method == http.MethodPost {
			// POST request
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddRole(item.Name)

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
		"roles/edit-permissions.html",
		map[string]interface{}{
			"Title": "Edit Role Permissions",
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func revokeRolePermissions(c echo.Context) error {
	var (
		name       = c.FormValue("name")
		method     = c.FormValue("_method")
		request    requests.RoleRevokePerm
		err        error
		etcdClient *v3.Client
		revokeList []types.KVPerm
	)

	if c.Request().Method == http.MethodPost && strings.EqualFold(method, http.MethodDelete) {
		etcdClient, err = v3.New(
			utils.GetEntryPoints(),
			utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
			utils.GetUsername(), utils.GetPassword(),
		)

		if err == nil {
			if bindErr := c.Bind(&request); bindErr == nil {
				for _, permHash := range request.PermHashes {
					permInput, inErr := requests.KVPerm{}.From(utils.Base64Decode(permHash))
					if inErr != nil {
						log.Errorf("revoke input error: %v", inErr)
						continue
					}

					revokeList = append(revokeList, types.KVPerm(permInput))
				}

				err = etcdClient.RevokePermissions(name, revokeList)

				if err != nil {
					return c.Render(
						http.StatusInternalServerError,
						"roles/edit-permissions.html",
						map[string]interface{}{
							"Title": "Edit Role Permissions",
							"Item":  types.RoleInfo{Name: name},
							"Error": utils.ErrorMessage(err),
							"csrf":  c.Get("csrf").(string),
						},
					)
				}
			}
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/role/edit/"+name)
}
