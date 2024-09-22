package home

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
)

func homeIndex(c echo.Context) error {
	return c.Render(
		http.StatusOK,
		"home/index.html",
		map[string]interface{}{
			"Title": "UI Configuration",
			"Settings": map[string]interface{}{
				"Endpoints":         utils.GetEntryPoints(),
				"Timeout":           utils.GetOpTimeout(),
				"UIWithCredentials": utils.GetUIUsername() != "",
				"WithCredentials":   utils.GetUsername() != "",
				"WithSSL":           utils.GetSSLCAFile() != "",
			},
		},
	)
}
