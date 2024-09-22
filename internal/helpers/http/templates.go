package http

import (
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/utking/etcd-ui/internal/controllers"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/templates"
	"github.com/utking/extemplate"
)

func InitTemplates(e *echo.Echo) error {
	xt := extemplate.New()

	if err := xt.ParseFS(templates.TemplateFiles, []string{".tmpl", ".html"}); err != nil {
		return err
	}

	t := &Template{
		worker: xt,
	}
	e.Renderer = t

	return nil
}

type Template struct {
	worker *extemplate.Extemplate
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	menu := controllers.GetMenu()

	return t.worker.ExecuteTemplate(w, name, map[string]interface{}{
		"data":    data,
		"title":   "Etcd UI",
		"year":    time.Now().Year(),
		"menu":    menu,
		"version": utils.GetReleaseVersion(),
	})
}
