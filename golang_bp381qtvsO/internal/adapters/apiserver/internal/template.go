package internal

import (
	"github.com/labstack/echo/v4"
	"io"
	"text/template"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RegisterTemplateRenderer(e *echo.Echo) {
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
}
