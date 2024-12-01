package pkg

import "fmt"

// processFile читает файл блоками и вызывает обработчик для каждого блока
func ProcessFile(data string, chunkSize int, handler func(string) (string, error)) (string, error) {
	finalText := ""
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		// Текущий блок строки
		chunk := data[i:end]
		// Вызываем обработчик с текущим блоком
		text, err := handler(chunk)
		fmt.Println(text)
		if err != nil {
			return "", err
		}
		finalText += text
	}

	return finalText, nil
}
