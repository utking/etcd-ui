package cluster

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/internal/providers/etcd/types"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
	"github.com/utking/etcd-ui/internal/types/requests"
)

func indexUsers(c echo.Context) error {
	var (
		code   = http.StatusOK
		users  = make([]types.UserRecord, 0)
		filter = c.QueryParam("filter")
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		users, err = etcdClient.GetUsers(filter)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"users/list.html",
		map[string]interface{}{
			"Title":  "Users",
			"Error":  utils.ErrorMessage(err),
			"Items":  users,
			"Filter": filter,
		},
	)
}

func infoUser(c echo.Context) error {
	var (
		code     = http.StatusOK
		username = c.QueryParam("name")
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

func createUser(c echo.Context) error {
	var (
		code       = http.StatusOK
		item       requests.UserCreate
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
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddUser(item.Name, item.Password)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/users?filter="+item.Name)
				}
			}
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"users/create.html",
		map[string]interface{}{
			"Title": "Create User",
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func passwdUser(c echo.Context) error {
	var (
		code       = http.StatusOK
		item       types.UserInfo
		request    requests.UserCreate
		err        error
		etcdClient *v3.Client
	)

	etcdClient, err = v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		if c.Request().Method == http.MethodGet {
			item, err = etcdClient.UserInfo(c.QueryParam("name"))
		} else if c.Request().Method == http.MethodPost {
			err = errors.Join(c.Bind(&request), request.Validate())
			item.Name = types.UserRecord(request.Name)

			if err == nil {
				err = etcdClient.ChangeUserPassword(request.Name, request.Password)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/users?filter="+request.Name)
				}
			}
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"users/passwd.html",
		map[string]interface{}{
			"Title": "Change User Password",
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func revokeUserGroups(c echo.Context) error {
	var (
		item       requests.UserEditRoles
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
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.RevokeUserRoles(item.Name, item.RevokeRoles)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/user?name="+item.Name)
				}
			}
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/user?name="+item.Name)
}

func addUserGroups(c echo.Context) error {
	var (
		code         = http.StatusOK
		item         requests.UserEditRoles
		err          error
		etcdClient   *v3.Client
		roles        []string
		allRoles     []string
		grantedRoles []string
	)

	etcdClient, err = v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		if c.Request().Method == http.MethodGet {
			item.Name = c.QueryParam("name")
			allRoles, _ = etcdClient.GetRoles("")

			if userInfo, userInfoErr := etcdClient.UserInfo(item.Name); userInfoErr == nil {
				grantedRoles = userInfo.Roles
				roles = make([]string, 0, max(len(allRoles)-len(grantedRoles), 0))

				for _, role := range allRoles {
					if !slices.Contains(grantedRoles, role) {
						roles = append(roles, role)
					}
				}
			}
		} else if c.Request().Method == http.MethodPost {
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddUserRoles(item.Name, item.AddRoles)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/user?name="+item.Name)
				}
			}
		}
	}

	return c.Render(
		code,
		"users/add-roles.html",
		map[string]interface{}{
			"Title":        "Grant Roles",
			"Item":         item,
			"Roles":        roles,
			"GrantedRoles": grantedRoles,
			"Error":        utils.ErrorMessage(err),
			"csrf":         c.Get("csrf").(string),
		},
	)
}

func deleteUser(c echo.Context) error {
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
			_ = etcdClient.DeleteUser(name)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/users?filter="+name)
}
