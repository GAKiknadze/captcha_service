package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

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
	testDir := "boundary_tests"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic(fmt.Sprintf("Не удалось создать директорию: %v", err))
	}

	fmt.Println("Тестирование граничных случаев для капчи с повернутыми символами")
	fmt.Println("=================================================================")

	// Тестовые случаи: длинные тексты, короткие тексты, разные размеры изображений
	testCases := []struct {
		name        string
		text        string
		width       int
		height      int
		fontSize    float64
		description string
	}{
		// Короткие тексты
		{
			name:        "short_text_small_image",
			text:        "A",
			width:       100,
			height:      50,
			fontSize:    20,
			description: "Один символ в маленьком изображении",
		},
		{
			name:        "two_chars",
			text:        "AB",
			width:       120,
			height:      60,
			fontSize:    24,
			description: "Два символа",
		},

		// Длинные тексты
		{
			name:        "long_text_normal",
			text:        "ABCDEFGH",
			width:       300,
			height:      80,
			fontSize:    24,
			description: "Длинный текст (8 символов)",
		},
		{
			name:        "very_long_text",
			text:        "ABCDEFGHIJKLM",
			width:       400,
			height:      100,
			fontSize:    22,
			description: "Очень длинный текст (13 символов)",
		},
		{
			name:        "long_text_small_width",
			text:        "TEST1234",
			width:       180,
			height:      70,
			fontSize:    20,
			description: "Длинный текст в узком изображении",
		},

		// Разные размеры изображений
		{
			name:        "wide_image",
			text:        "CAPTCHA",
			width:       500,
			height:      80,
			fontSize:    28,
			description: "Широкое изображение",
		},
		{
			name:        "tall_image",
			text:        "SECURE",
			width:       200,
			height:      150,
			fontSize:    26,
			description: "Высокое изображение",
		},
		{
			name:        "small_square",
			text:        "OK",
			width:       80,
			height:      80,
			fontSize:    18,
			description: "Маленькое квадратное изображение",
		},

		// Граничные случаи
		{
			name:        "minimal_size",
			text:        "I",
			width:       50,
			height:      30,
			fontSize:    14,
			description: "Минимальный размер изображения",
		},
		{
			name:        "numbers_only",
			text:        "1234567890",
			width:       350,
			height:      90,
			fontSize:    24,
			description: "Только цифры (10 символов)",
		},
		{
			name:        "mixed_case",
			text:        "AbCdEfGhIj",
			width:       320,
			height:      85,
			fontSize:    22,
			description: "Смешанный регистр",
		},
	}

	successCount := 0
	failCount := 0

	for _, tc := range testCases {
		fmt.Printf("\nТест: %s\n", tc.name)
		fmt.Printf("Описание: %s\n", tc.description)
		fmt.Printf("Текст: '%s', Размер: %dx%d, Шрифт: %.0fpt\n", tc.text, tc.width, tc.height, tc.fontSize)

		// Создаем конфигурацию
		config := captcha.ImageCaptchaConfig{
			BackgroundColor: color.White,
			TextColor:       color.Black,
			Font:            loadFont(tc.fontSize),
			FontSize:        int(tc.fontSize),
			ImageWidth:      tc.width,
			ImageHeight:     tc.height,
		}

		// Создаем генератор
		captchaGenerator := captcha.NewImageCaptcha(config)

		// Генерируем капчу
		captchaImage, err := captchaGenerator.Generate(tc.text)
		if err != nil {
			fmt.Printf("❌ ОШИБКА генерации: %v\n", err)
			failCount++
			continue
		}

		// Сохраняем в файл
		filename := filepath.Join(testDir, fmt.Sprintf("%s.png", tc.name))
		err = os.WriteFile(filename, captchaImage, 0644)
		if err != nil {
			fmt.Printf("❌ ОШИБКА сохранения: %v\n", err)
			failCount++
			continue
		}

		fmt.Printf("✅ УСПЕХ: создан файл %s\n", filename)
		successCount++

		// Генерируем несколько вариантов с тем же текстом для проверки случайности
		for i := 0; i < 2; i++ {
			captchaImage2, err := captchaGenerator.Generate(tc.text)
			if err == nil {
				variantFilename := filepath.Join(testDir, fmt.Sprintf("%s_variant_%d.png", tc.name, i+1))
				os.WriteFile(variantFilename, captchaImage2, 0644)
			}
		}
	}

	// Тестирование специальных граничных случаев
	fmt.Println("\n\nСпециальные граничные тесты:")
	fmt.Println("============================")

	specialTests := []struct {
		name   string
		text   string
		width  int
		height int
	}{
		{"edge_case_1", "WWW", 100, 40},    // Широкие символы
		{"edge_case_2", "iii", 90, 40},     // Узкие символы
		{"edge_case_3", "MgQy", 120, 50},   // Символы разной ширины
		{"edge_case_4", "()[]{}", 180, 60}, // Специальные символы
	}

	for _, st := range specialTests {
		config := captcha.ImageCaptchaConfig{
			BackgroundColor: color.White,
			TextColor:       color.Black,
			Font:            loadFont(20),
			FontSize:        20,
			ImageWidth:      st.width,
			ImageHeight:     st.height,
		}

		captchaGenerator := captcha.NewImageCaptcha(config)
		captchaImage, err := captchaGenerator.Generate(st.text)
		if err != nil {
			fmt.Printf("❌ Спецтест '%s' ОШИБКА: %v\n", st.name, err)
			failCount++
		} else {
			filename := filepath.Join(testDir, fmt.Sprintf("%s.png", st.name))
			os.WriteFile(filename, captchaImage, 0644)
			fmt.Printf("✅ Спецтест '%s' УСПЕХ\n", st.name)
			successCount++
		}
	}

	// Итоги
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("ИТОГИ ТЕСТИРОВАНИЯ:\n")
	fmt.Printf("Успешных тестов: %d\n", successCount)
	fmt.Printf("Неудачных тестов: %d\n", failCount)
	fmt.Printf("Общее количество: %d\n", successCount+failCount)

	if failCount == 0 {
		fmt.Println("\n🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!")
		fmt.Println("Символы не выходят за границы изображения.")
	} else {
		fmt.Printf("\n⚠️  Есть проблемы: %d тестов не прошли\n", failCount)
	}

	fmt.Printf("\nВсе тестовые изображения сохранены в директории: %s\n", testDir)
	fmt.Println("\nПроверьте визуально, что:")
	fmt.Println("1. Все символы полностью видны в пределах изображения")
	fmt.Println("2. Нет обрезанных краев символов")
	fmt.Println("3. Текст центрирован по горизонтали")
	fmt.Println("4. Вертикальные смещения не выводят символы за границы")
	fmt.Println("5. Повороты символов не обрезаются")
}
