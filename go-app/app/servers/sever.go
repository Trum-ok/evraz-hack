package servers

import (
	"github.com/Danila331/hach-evroasia/app/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() error {
	app := echo.New()
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(500, map[string]string{"error": err.Error()})
	}
	app.Static("/downloads", "./uploads")
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.GET("/favicon.ico", func(c echo.Context) error {
		return nil // Убедитесь, что путь правильный
	})
	app.GET("/", handlers.UploadFormHandler)
	app.POST("/", handlers.UploadFileHandler)
	app.GET("/download/:file", handlers.DownloadFileHandler)
	return app.Start(":80")
}
