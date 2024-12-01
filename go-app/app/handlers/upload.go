package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Danila331/hach-evroasia/app/pkg"
	"github.com/labstack/echo/v4"
)

var Count int = 0

func UploadFormHandler(c echo.Context) error {
	// Пути к HTML-шаблону

	htmlFiles := []string{
		filepath.Join("./", "templates", "upload.html"),
	}

	// Рендеринг HTML-формы
	templ, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		return err
	}

	return templ.ExecuteTemplate(c.Response(), "upload", nil)
}

func UploadFileHandler(c echo.Context) error {
	// Формируем директорию для сохранения файлов
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.Logger().Error("Произошла ошибка: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось создать директорию для загрузки",
		})
	}

	// Перебираем все файлы из запроса
	form, err := c.MultipartForm()
	if err != nil {
		c.Logger().Error("Произошла ошибка: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Ошибка при чтении формы",
		})
	}
	files := form.File["file"]
	finalText := ""
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		if extension == ".zip" {
			// Обработка zip файлов
			_, err = pkg.OpenAllCFilesInZipWithStructure(file)
			if err != nil {
				c.Logger().Error("Произошла ошибка: ", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": fmt.Sprintf("Ошибка при обработке zip файла: %v", err),
				})
			}
			// Возможно вы хотите обработать структуру файлов или вернуть её клиенту
			// Формируем URL для скачивания .md файла
			mdFileURL := fmt.Sprintf("/downloads/%s", filepath.Base("output.md"))

			// Возвращаем JSON с результатом
			return c.JSON(http.StatusOK, map[string]interface{}{
				"Message":   "Code-reiwew проведен успешно, ниже кнопка для скачивания .md файла",
				"MDFileURL": mdFileURL,
			})
			// Тут выводим структуру
		} else {
			// Открываем файл из формы
			src, err := file.Open()
			if err != nil {
				c.Logger().Error("Произошла ошибка: ", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Ошибка при открытии файла",
				})
			}
			defer src.Close()
			content, err := io.ReadAll(src)
			if err != nil {
				c.Logger().Error("Произошла ошибка: ", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Ошибка при чтении файла",
				})
			}
			fmt.Println(len(string(content)))
			text := ""
			if extension == ".cs" {
				text, err = pkg.ProcessFile(string(content), 5200, pkg.СGetRequestToLlm)
				if err != nil {
					c.Logger().Error("Произошла ошибка: ", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": fmt.Sprintf("Ошибка при обработке файла: %v", err),
					})
				}
			}
			if extension == ".py" {
				text, err = pkg.ProcessFile(string(content), 5000, pkg.PythonGetRequestToLlm)
				if err != nil {
					c.Logger().Error("Произошла ошибка: ", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": fmt.Sprintf("Ошибка при обработке файла: %v", err),
					})
				}
			}
			if extension == ".js" || extension == ".ts" || extension == ".tsx" {
				text, err = pkg.ProcessFile(string(content), 5200, pkg.JsGetRequestToLlm)
				if err != nil {
					c.Logger().Error("Произошла ошибка: ", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": fmt.Sprintf("Ошибка при обработке файла: %v", err),
					})
				}
			}
			if extension == ".go" {
				text, err = pkg.ProcessFile(string(content), 5200, pkg.JsGetRequestToLlm)
				if err != nil {
					c.Logger().Error("Произошла ошибка: ", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": fmt.Sprintf("Ошибка при обработке файла: %v", err),
					})
				}
			}
			finalText += text + "\n"
		}
	}

	// Генерация .md файла
	mdFileName := "files.md"
	mdFile, err := os.Create("./uploads/" + mdFileName)
	if err != nil {
		c.Logger().Error("Произошла ошибка: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Ошибка при создании .md файла: %v", err),
		})
	}
	defer mdFile.Close()

	// Записываем ссылки на сохраненные файлы в .md файл
	mdFile.WriteString(finalText)

	// Формируем URL для скачивания .md файла
	mdFileURL := fmt.Sprintf("/downloads/%s", filepath.Base(mdFileName))

	// Возвращаем JSON с результатом
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Message":   "Code-reiwew проведен успешно, ниже кнопка для скачивания .md файла",
		"MDFileURL": mdFileURL,
	})
}
