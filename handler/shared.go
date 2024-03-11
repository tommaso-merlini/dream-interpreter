package handler

import (
	"log/slog"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Make(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := h(c)
		if err != nil {
			slog.Error("internal server error", "err", err, "path", c.Request().URL.Path)
		}
		return err
	}
}

func render(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/html")
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func hxRedirect(c echo.Context, to string) {
	c.Response().Header().Set("HX-Redirect", to)
}
