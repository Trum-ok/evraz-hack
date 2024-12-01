package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Список допустимых расширений файлов
var allowedExtensions = []string{".go", ".cs", ".py", ".js"}

// Функция для проверки, является ли файл поддерживаемым
func isAllowedExtension(filename string) bool {
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func OpenAllCFilesInZipWithStructure(fileHeader *multipart.FileHeader) (string, error) {
	// Открываем загруженный файл
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть загруженный файл: %v", err)
	}
	defer src.Close()

	// Создаем временный файл для работы с архивом
	tempFile, err := os.CreateTemp("", "uploaded-*.zip")
	if err != nil {
		return "", fmt.Errorf("не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Сохраняем содержимое в временный файл
	if _, err := io.Copy(tempFile, src); err != nil {
		return "", fmt.Errorf("не удалось сохранить содержимое файла: %v", err)
	}

	// Открываем zip-архив из временного файла
	r, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("не удалось открыть архив: %v", err)
	}
	defer r.Close()

	// Карта для содержимого файлов
	filesContent := make(map[string]string)
	// Буфер для формирования файловой структуры
	var treeBuffer bytes.Buffer
	treeBuffer.WriteString("Archive structure:\n")

	// Карта для отслеживания добавленных папок
	seenDirs := make(map[string]bool)
	outputFile, err := os.Create("./uploads/output.md")
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл для записи: %v", err)
	}
	defer outputFile.Close()

	// Mutex для безопасной записи в файл
	var mu sync.Mutex
	// WaitGroup для управления go-routines
	var wg sync.WaitGroup

	// Канал для ограничения количества одновременно работающих горутин
	concurrentLimit := 4
	sem := make(chan struct{}, concurrentLimit)

	// Читаем файлы из архива
	for _, file := range r.File {
		// Определяем уровень файла/директории
		parts := strings.Split(file.Name, "/")
		level := len(parts) - 1

		// Добавляем недостающие папки в структуру дерева
		for i := 1; i <= level; i++ {
			dirPath := strings.Join(parts[:i], "/")
			if !seenDirs[dirPath] {
				indent := strings.Repeat(" │   ", i-1)
				treeBuffer.WriteString(fmt.Sprintf("%s ├── %s/\n", indent, parts[i-1]))
				seenDirs[dirPath] = true
			}
		}

		// Обрабатываем файлы
		if !file.FileInfo().IsDir() {
			indent := strings.Repeat(" │   ", level-1)
			treeBuffer.WriteString(fmt.Sprintf("%s ├── %s\n", indent, parts[level]))

			// Если файл разрешённого расширения
			if isAllowedExtension(file.Name) {
				wg.Add(1)

				// Ограничиваем количество одновременно работающих горутин
				sem <- struct{}{} // Запускаем горутину, добавляем элемент в канал
				go func(file *zip.File) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Паника в goroutine: %v\n", r)
						}
					}()
					// Убираем элемент из канала по завершении горутины
					defer func() { <-sem }() // Имитация задержки при обработке файла

					// Читаем содержимое файла
					srcFile, err := file.Open()
					if err != nil {
						fmt.Printf("не удалось открыть файл внутри архива: %v\n", err)
						return
					}
					defer srcFile.Close()

					content, err := io.ReadAll(srcFile)
					if err != nil {
						fmt.Printf("не удалось прочитать содержимое файла: %v\n", err)
						return
					}

					// Отправляем содержимое на LLM
					ext := filepath.Ext(file.Name)
					text := ""
					if ext == ".cs" {
						if len(content) > 5200 {
							content = content[:5200]
						}
						text, err = СGetRequestToLlm(string(content))
						if err != nil {
							fmt.Printf("не удалось отправить запрос на LLM: %v\n", err)
							return
						}
					}
					if ext == ".py" {
						if len(content) > 5200 {
							content = content[:5200]
						}
						text, err = PythonGetRequestToLlm(string(content))
						if err != nil {
							fmt.Printf("не удалось отправить запрос на LLM: %v\n", err)
							return
						}
					}

					if ext == ".js" || ext == ".ts" || ext == ".tsx" {
						if len(content) > 5400 {
							content = content[:5200]
						}
						text, err = JsGetRequestToLlm(string(content))
						if err != nil {
							fmt.Printf("не удалось отправить запрос на LLM: %v\n", err)
							return
						}
					}
					if ext == ".go" {
						if len(content) > 5400 {
							content = content[:5200]
						}
						text, err = GoGetRequestToLlm(string(content))
						if err != nil {
							fmt.Printf("не удалось отправить запрос на LLM: %v\n", err)
							return
						}
					}

					// Безопасно записываем в выходной файл
					mu.Lock()
					defer mu.Unlock()
					_, err = outputFile.WriteString(fmt.Sprintf("File: %s\n%s\n\n", file.Name, text))
					if err != nil {
						fmt.Printf("не удалось записать в файл: %v\n", err)
					}
				}(file) // Передаём текущий файл как аргумент в go-routine
			} else {
				// Если файл не разрешён, добавляем пустое содержимое
				filesContent[file.Name] = "Файл не читается"
			}

		}
	}

	// Ожидаем завершения всех потоков
	wg.Wait()
	structure := treeBuffer.String()
	text, err := CArchitectureGetRequestToLlm(structure)
	if err != nil {
		fmt.Printf("Ошибка в ArchitectureGetRequestToLlm: %v\n", err)
	}
	_, err = outputFile.WriteString(text)
	if err != nil {
		fmt.Printf("не удалось записать в файл: %v\n", err)
	}
	// Возвращаем файловую структуру
	return treeBuffer.String(), nil
}
