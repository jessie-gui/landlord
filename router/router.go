package router

import (
	"github/jessie-gui/landlord/controller"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Template /**
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewEcho() *echo.Echo {
	e := echo.New()

	e.Static("/static", "static")
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.GET("/", controller.Index)
	e.GET("/open", controller.Connect)

	return e
}
