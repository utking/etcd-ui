package home

import (
	"github.com/labstack/echo/v4"
	types "github.com/utking/etcd-ui/internal/types/common"
)

func Setup(app *echo.Echo, _ *types.WebMenu) {
	app.GET("/", homeIndex)
}
