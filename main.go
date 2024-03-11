package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/tommaso-merlini/dream-interpreter/handler"
)

//go:embed public
var FS embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e := echo.New()
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(FS)))))
	e.Static("/images", "./images")
	e.GET("/", handler.Make(handler.HomeShow))
	e.GET("/chat", handler.Make(handler.ChatShow))
	e.GET("/chatws", handler.Make(handler.ChatWS))
	e.DELETE("/thinking-message", handler.Make(handler.DeleteThinkingMessage))
	e.POST("/add-email", handler.Make(handler.AddEmail))

	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}
