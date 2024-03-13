package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/tommaso-merlini/dream-interpreter/view/home"
)

func HomeShow(c echo.Context) error {
	return render(c, home.Home())
}

func CopyToClipboard(c echo.Context) error {
	return c.String(200, "Copied!")
}
