package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/labstack/echo/v4"
)

type CustomErrorData struct {
	StatusCode int
	Message    string
}

func CustomErrorHandler(err error, c echo.Context) {
	// Определяем HTTP-статус и сообщение
	c.Logger().Error("Произошла ошибка: ", err)
	var statusCode int
	var message string

	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode = httpErr.Code
		if httpErr.Message != nil {
			message = httpErr.Message.(string)
		} else {
			message = http.StatusText(httpErr.Code)
		}
	} else {
		statusCode = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)
	}

	// Пути к HTML-шаблону
	htmlFiles := []string{
		filepath.Join("./", "templates", "error.html"),
	}

	// Рендеринг HTML-страницы
	templ, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		c.Logger().Error("Ошибка при загрузке HTML-шаблона: ", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Передаем данные в шаблон
	data := CustomErrorData{
		StatusCode: statusCode,
		Message:    message,
	}

	// Устанавливаем статус и рендерим
	c.Response().WriteHeader(statusCode)
	if err := templ.ExecuteTemplate(c.Response(), "error", data); err != nil {
		c.Logger().Error("Ошибка при рендеринге HTML-шаблона: ", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
