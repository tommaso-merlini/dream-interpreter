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
	// e.Use(handler.RateLimiterMiddleware(100, 24*time.Hour))

	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(FS)))))
	e.Static("/images", "./images")
	e.GET("/emails", handler.Make(handler.EmailDump))
	e.GET("/sitemap.xml", func(c echo.Context) error {
		return c.File("./sitemap.xml")
	})
	e.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./robots.txt")
	})

	e.GET("/", handler.Make(handler.HomeShow))
	e.GET("/chat", handler.Make(handler.ChatShow))
	e.GET("/chatws", handler.Make(handler.ChatWS))
	e.DELETE("/thinking-message", handler.Make(handler.DeleteThinkingMessage))
	e.POST("/add-email", handler.Make(handler.AddEmail))
	e.GET("/copy-to-clipboard", handler.Make(handler.CopyToClipboard))

	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}
