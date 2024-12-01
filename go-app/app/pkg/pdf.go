package pkg

import (
	"log"

	"github.com/jung-kurt/gofpdf"
)

// CreatePDF создает PDF-файл с указанным текстом
func CreatePDF(filename string, text string) error {
	// Создаем новый PDF документ
	pdf := gofpdf.New("P", "mm", "A4", "") // P - портретная ориентация, mm - миллиметры, A4 - размер страницы

	pdf.AddUTF8Font("DejaVuSans", "", "./static/fonts/DejaVuSans.ttf") // Добавляем шрифт

	// Добавляем страницу
	pdf.SetFont("DejaVuSans", "", 12)
	// Добавляем текст
	pdf.Cell(40, 10, "Заголовок:")
	pdf.Ln(10) // Добавляет разрыв строки

	// Добавляем основной текст с многострочным выводом
	pdf.MultiCell(0.0, 10.0, text, "", "L", false)

	// Сохраняем PDF в файл
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		log.Println("Ошибка при создании PDF:", err)
		return err
	}

	log.Println("PDF файл создан:", filename)
	return nil
}
