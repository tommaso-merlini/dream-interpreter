package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tommaso-merlini/dream-interpreter/handler"
)

//go:embed public
var FS embed.FS

func main() {
	e := echo.New()
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(FS)))))
	e.Static("/images", "./images")
	e.GET("/", handler.Make(handler.HomeShow))
	e.GET("/chat", handler.Make(handler.ChatShow))
	e.GET("/chatws", handler.Make(handler.ChatWS))
	e.DELETE("/thinking-message", handler.Make(handler.DeleteThinkingMessage))

	e.Logger.Fatal(e.Start(":3000"))
}
