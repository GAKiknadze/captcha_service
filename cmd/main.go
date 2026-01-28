package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/GAKiknadze/captcha_service/internal/captcha"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// loadFont загружает и создает шрифт с указанным размером
func loadFont(size float64) *font.Face {
	// Парсим шрифт Go Regular
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	// Создаем шрифт с указанным размером
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return &face
}

func main() {
	// Инициализация конфигурации капчи
	config := captcha.ImageCaptchaConfig{
		BackgroundColor: color.White,
		TextColor:       color.Black,
		Font:            loadFont(28), // Загружаем шрифт Go Regular с размером 28
		FontSize:        28,
		ImageWidth:      250,
		ImageHeight:     100,
	}

	// Создание генератора капчи
	captchaGenerator := captcha.NewImageCaptcha(config)

	// Тестирование разных текстов для капчи
	testTexts := []string{
		"ABCD123",
		"CAPTCHA",
		"SECURE42",
		"TEST789",
		"RANDOM",
		"123456",
		"VERIFY",
		"ACCESS",
	}

	// Генерируем капчи для каждого текста
	for i, captchaText := range testTexts {
		captchaImage, err := captchaGenerator.Generate(captchaText)
		if err != nil {
			fmt.Printf("Ошибка генерации капчи '%s': %v\n", captchaText, err)
			continue
		}

		// Сохранение капчи в файл
		filename := fmt.Sprintf("captcha_%d_%s.png", i, captchaText)
		err = os.WriteFile(filename, captchaImage, 0644)
		if err != nil {
			fmt.Printf("Ошибка сохранения файла %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Капча '%s' успешно сохранена в файл %s\n", captchaText, filename)
	}

	// Генерация нескольких капч с одинаковым текстом для демонстрации случайности
	fmt.Println("\nГенерация 3 капч с одинаковым текстом 'TEST123':")
	for i := 0; i < 3; i++ {
		captchaImage, err := captchaGenerator.Generate("TEST123")
		if err != nil {
			fmt.Printf("Ошибка генерации капчи %d: %v\n", i, err)
			continue
		}

		filename := fmt.Sprintf("same_text_%d_TEST123.png", i)
		err = os.WriteFile(filename, captchaImage, 0644)
		if err != nil {
			fmt.Printf("Ошибка сохранения файла %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Капча %d сохранена в файл %s\n", i, filename)
	}

	fmt.Println("\nВсе капчи успешно сгенерированы!")
	fmt.Println("Особенности реализованного поворота символов:")
	fmt.Println("1. Каждый символ поворачивается на случайный угол от -20° до +20°")
	fmt.Println("2. Добавлено случайное вертикальное смещение символов")
	fmt.Println("3. Текст автоматически центрируется с учетом поворотов")
	fmt.Println("4. Символы не выходят за границы изображения")
	fmt.Println("5. Применено волнообразное искажение изображения")
}
