package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

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
	// Создаем папку для тестовых капч
	testDir := "test_captchas"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic(fmt.Sprintf("Не удалось создать директорию: %v", err))
	}

	// Тестовые тексты для капчи
	testTexts := []string{
		"ABCD123",
		"XYZ789",
		"CAPTCHA",
		"SECURE",
		"RANDOM",
		"123456",
		"TESTME",
		"VERIFY",
		"ACCESS",
		"LOGIN",
	}

	// Конфигурации для тестирования
	configs := []struct {
		name   string
		config captcha.ImageCaptchaConfig
	}{
		{
			name: "basic",
			config: captcha.ImageCaptchaConfig{
				BackgroundColor: color.White,
				TextColor:       color.Black,
				Font:            loadFont(24),
				FontSize:        24,
				ImageWidth:      200,
				ImageHeight:     80,
			},
		},
		{
			name: "colored",
			config: captcha.ImageCaptchaConfig{
				BackgroundColor: color.RGBA{R: 240, G: 240, B: 240, A: 255},
				TextColor:       color.RGBA{R: 50, G: 100, B: 200, A: 255},
				Font:            loadFont(28),
				FontSize:        28,
				ImageWidth:      220,
				ImageHeight:     90,
			},
		},
		{
			name: "dark",
			config: captcha.ImageCaptchaConfig{
				BackgroundColor: color.RGBA{R: 30, G: 30, B: 30, A: 255},
				TextColor:       color.RGBA{R: 220, G: 220, B: 100, A: 255},
				Font:            loadFont(26),
				FontSize:        26,
				ImageWidth:      250,
				ImageHeight:     100,
			},
		},
	}

	// Генерируем капчи для каждой конфигурации и каждого текста
	for _, cfg := range configs {
		captchaGenerator := captcha.NewImageCaptcha(cfg.config)

		for i, text := range testTexts {
			// Генерируем капчу
			captchaImage, err := captchaGenerator.Generate(text)
			if err != nil {
				fmt.Printf("Ошибка генерации капчи %s-%d: %v\n", cfg.name, i, err)
				continue
			}

			// Сохраняем в файл
			filename := filepath.Join(testDir, fmt.Sprintf("captcha_%s_%d_%s.png", cfg.name, i, text))
			err = os.WriteFile(filename, captchaImage, 0644)
			if err != nil {
				fmt.Printf("Ошибка сохранения файла %s: %v\n", filename, err)
				continue
			}

			fmt.Printf("Создана капча: %s (текст: %s)\n", filename, text)
		}
	}

	// Создаем несколько капч с одинаковым текстом для демонстрации случайности
	fmt.Println("\nГенерация нескольких капч с одинаковым текстом для демонстрации случайности:")
	sameText := "TEST123"
	baseConfig := captcha.ImageCaptchaConfig{
		BackgroundColor: color.White,
		TextColor:       color.Black,
		Font:            loadFont(24),
		FontSize:        24,
		ImageWidth:      200,
		ImageHeight:     80,
	}

	for i := 0; i < 5; i++ {
		captchaGenerator := captcha.NewImageCaptcha(baseConfig)
		captchaImage, err := captchaGenerator.Generate(sameText)
		if err != nil {
			fmt.Printf("Ошибка генерации капчи %d: %v\n", i, err)
			continue
		}

		filename := filepath.Join(testDir, fmt.Sprintf("same_text_%d_%s.png", i, sameText))
		err = os.WriteFile(filename, captchaImage, 0644)
		if err != nil {
			fmt.Printf("Ошибка сохранения файла %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Создана капча %d с текстом '%s': %s\n", i, sameText, filename)
	}

	fmt.Printf("\nВсе капчи сохранены в директории: %s\n", testDir)
	fmt.Println("Особенности реализованного поворота символов:")
	fmt.Println("1. Каждый символ поворачивается на случайный угол от -20° до +20°")
	fmt.Println("2. Добавлено случайное вертикальное смещение символов")
	fmt.Println("3. Межсимвольные интервалы слегка варьируются")
	fmt.Println("4. Применено волнообразное искажение изображения")
	fmt.Println("5. Добавлены случайные помехи (точки и линии)")
}
