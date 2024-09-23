package home

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	v3 "github.com/utking/etcd-ui/internal/providers/etcd/v3"
)

func homeIndex(c echo.Context) error {
	authEnabled, authCheckErr := v3.CheckAuthEnabled(
		utils.GetEntryPoints(),
		utils.GetSSLCertFile(), utils.GetSSLKeyFile(), utils.GetSSLCAFile(),
	)

	if authCheckErr != nil {
		log.Printf("error checking if auth is enabled: %+v\n", authCheckErr)
	}

	return c.Render(
		http.StatusOK,
		"home/index.html",
		map[string]interface{}{
			"Title": "UI Configuration",
			"Settings": map[string]interface{}{
				"Endpoints":         utils.GetEntryPoints(),
				"Timeout":           utils.GetOpTimeout(),
				"UIWithCredentials": utils.GetUIUsername() != "",
				"AuthEnabled":       authEnabled,
				"WithCredentials":   utils.GetUsername() != "",
				"WithSSL":           utils.GetSSLCAFile() != "",
			},
		},
	)
}
