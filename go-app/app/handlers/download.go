package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func DownloadFileHandler(c echo.Context) error {
	filePath := c.Param("filePath")

	// Открываем файл для скачивания
	file, err := os.Open(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Ошибка при открытии файла: %v", err),
		})
	}
	defer file.Close()

	// Устанавливаем заголовки для скачивания
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	c.Response().Header().Set("Content-Type", "application/octet-stream")

	// Отправляем файл клиенту
	_, err = io.Copy(c.Response(), file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Ошибка при отправке файла: %v", err),
		})
	}

	return nil
}
