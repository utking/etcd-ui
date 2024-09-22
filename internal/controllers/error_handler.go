package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/helpers/utils"
)

func HTTPErrorHandler(err error, c echo.Context) {
	templateFile := "errors/error.html"

	code := http.StatusInternalServerError

	if c.Response().Committed {
		return
	}

	// Check if there is an existing code to use
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		code = httpError.Code
	}

	if httpError != nil {
		switch httpError.Code {
		case http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError:
			templateFile = fmt.Sprintf("errors/%d.html", code)
		case http.StatusMethodNotAllowed:
			_ = c.JSON(
				http.StatusMethodNotAllowed,
				map[string]interface{}{
					"Error": utils.ErrorMessage(err),
				})

			return
		default:
			templateFile = "errors/error.html"
		}
	}

	errPageRenderError := c.Render(
		code,
		templateFile,
		map[string]interface{}{
			"Title": fmt.Sprintf("Error %d", code),
			"Error": err,
		})

	// In case the nice HTML render fails
	if errPageRenderError != nil {
		_ = c.JSON(code, map[string]interface{}{
			"Title": fmt.Sprintf("Error %d", code),
			"Error": err,
		})
	}
}
