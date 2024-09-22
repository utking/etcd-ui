package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/controllers/cluster"
	"github.com/utking/etcd-ui/internal/controllers/home"
	types "github.com/utking/etcd-ui/internal/types/common"
)

var (
	webMenu *types.WebMenu
)

func Setup(app *echo.Echo) {
	webMenu = new(types.WebMenu)
	webMenu.AdminItems = make([]map[string][]types.WebMenuItem, 0)
	webMenu.UserItems = make([]map[string][]types.WebMenuItem, 0)

	home.Setup(app, webMenu)
	cluster.Setup(app, webMenu)
}

func GetMenu() types.WebMenu {
	if webMenu == nil {
		return types.WebMenu{}
	}

	return *webMenu
}
