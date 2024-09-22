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

const (
	keysListLimit = 250
)

func listKeys(c echo.Context) error {
	var (
		code   = http.StatusOK
		filter = c.QueryParam("filter")
		items  []types.KVItem
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)
	if err == nil {
		items, err = etcdClient.ListKeys(
			filter,        // prefix to search by
			keysListLimit, // limit to get
			false,         // without data
		)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"keys/list.html",
		map[string]interface{}{
			"Title":  "List keys",
			"Error":  utils.ErrorMessage(err),
			"Items":  items,
			"Filter": filter,
			"Limit":  keysListLimit,
		},
	)
}

func infoKey(c echo.Context) error {
	var (
		code = http.StatusOK
		key  = c.QueryParam("key")
		item *types.KVItem
	)

	etcdClient, err := v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		item, err = etcdClient.Get(key)
	}

	if err != nil {
		code = http.StatusNotFound
	}

	return c.Render(
		code,
		"keys/info.html",
		map[string]interface{}{
			"Title":  "Key Info",
			"Error":  utils.ErrorMessage(err),
			"KVItem": item,
			"csrf":   c.Get("csrf").(string),
		},
	)
}

func createKey(c echo.Context) error {
	var (
		code       = http.StatusOK
		item       requests.KVCreate
		err        error
		etcdClient *v3.Client
	)

	etcdClient, err = v3.New(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
		utils.GetUsername(), utils.GetPassword(),
	)

	if err == nil {
		item.Key = c.QueryParam("key")

		if c.Request().Method == http.MethodGet && item.Key != "" {
			kvItem, getErr := etcdClient.Get(item.Key)
			if getErr == nil {
				item.TTL = kvItem.LeaseTTL
				item.Value = kvItem.Value.Unwrap()
			}
		} else if c.Request().Method == http.MethodPost {
			err = c.Bind(&item)
			if err == nil {
				err = etcdClient.Put(item.Key, item.Value, item.TTL)

				if err == nil {
					return c.Redirect(http.StatusSeeOther, "/cluster/keys/list")
				}
			}
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	return c.Render(
		code,
		"keys/create.html",
		map[string]interface{}{
			"Title": "Put Key",
			"Item":  item,
			"Error": utils.ErrorMessage(err),
			"csrf":  c.Get("csrf").(string),
		},
	)
}

func deleteKey(c echo.Context) error {
	var (
		key        = c.FormValue("key")
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
			_ = etcdClient.Delete(key)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/cluster/keys/list")
}
