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
		// GET request
		item.Name = c.Param("name")
		if c.Request().Method == http.MethodGet && item.Name != "" {
			userItem, getErr := etcdClient.UserInfo(item.Name)
			if getErr == nil {
				item.Name = string(userItem.Name)
				item.Roles = userItem.Roles
			}
		} else if c.Request().Method == http.MethodPost {
			// POST request
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.AddUser(item.Name, item.Password)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/users/list?filter="+item.Name)
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
			"Title": "Create User", // FIX: for edit
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
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

	return c.Redirect(http.StatusSeeOther, "/cluster/users/list?filter="+name)
}
