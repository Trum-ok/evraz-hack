package servers

import (
	"net/http"

	"github.com/Danila331/hach-evroasia/app/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() error {
	app := echo.New()
	app.HTTPErrorHandler = handlers.CustomErrorHandler
	app.Static("/downloads", "./uploads")
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.GET("/favicon.ico", func(c echo.Context) error {
		return nil // Убедитесь, что путь правильный
	})
	app.GET("/test-error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "Страница не найдена")
	})
	app.GET("/", handlers.UploadFormHandler)
	app.POST("/", handlers.UploadFileHandler)
	app.GET("/download/:file", handlers.DownloadFileHandler)
	return app.Start(":80")
}
