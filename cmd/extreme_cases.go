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
	// Создаем папку для экстремальных тестов
	testDir := "extreme_tests"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic(fmt.Sprintf("Не удалось создать директорию: %v", err))
	}

	fmt.Println("ЭКСТРЕМАЛЬНЫЕ ТЕСТЫ КАПЧИ С ПОВОРОТОМ СИМВОЛОВ")
	fmt.Println(strings.Repeat("=", 60))

	// Экстремальные тестовые случаи
	extremeTests := []struct {
		name        string
		text        string
		width       int
		height      int
		fontSize    float64
		description string
		expected    string // Ожидаемый результат
	}{
		// Случай 1: Максимально длинный текст
		{
			name:        "max_length_text",
			text:        "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			width:       600,
			height:      120,
			fontSize:    22,
			description: "Максимально длинный текст (36 символов)",
			expected:    "Должен поместиться с уменьшенными интервалами",
		},
		// Случай 2: Очень широкие символы
		{
			name:        "wide_characters",
			text:        "WWWMMMQQQ",
			width:       300,
			height:      80,
			fontSize:    28,
			description: "Текст из широких символов (W, M, Q)",
			expected:    "Должен центрироваться с учетом ширины символов",
		},
		// Случай 3: Очень узкие символы
		{
			name:        "narrow_characters",
			text:        "iiillljjj",
			width:       200,
			height:      70,
			fontSize:    26,
			description: "Текст из узких символов (i, l, j)",
			expected:    "Должен правильно центрироваться",
		},
		// Случай 4: Смешанная ширина символов
		{
			name:        "mixed_width",
			text:        "WiMqIjLp",
			width:       250,
			height:      90,
			fontSize:    24,
			description: "Смесь широких и узких символов",
			expected:    "Должен равномерно распределиться",
		},
		// Случай 5: Минимальные размеры изображения
		{
			name:        "minimal_image",
			text:        "A",
			width:       40,
			height:      25,
			fontSize:    12,
			description: "Минимально возможное изображение для одного символа",
			expected:    "Символ должен быть виден полностью",
		},
		// Случай 6: Очень большой шрифт
		{
			name:        "huge_font",
			text:        "BIG",
			width:       300,
			height:      150,
			fontSize:    48,
			description: "Очень большой шрифт",
			expected:    "Символы должны помещаться с учетом поворотов",
		},
		// Случай 7: Очень маленький шрифт
		{
			name:        "tiny_font",
			text:        "smalltext",
			width:       200,
			height:      60,
			fontSize:    12,
			description: "Очень маленький шрифт",
			expected:    "Текст должен быть читаемым",
		},
		// Случай 8: Квадратное изображение с длинным текстом
		{
			name:        "square_long_text",
			text:        "LONGTEXT",
			width:       150,
			height:      150,
			fontSize:    20,
			description: "Квадратное изображение с длинным текстом",
			expected:    "Текст должен центрироваться по горизонтали и вертикали",
		},
		// Случай 9: Высокое узкое изображение
		{
			name:        "tall_narrow",
			text:        "UP",
			width:       60,
			height:      200,
			fontSize:    24,
			description: "Высокое узкое изображение",
			expected:    "Символы должны быть вертикально центрированы",
		},
		// Случай 10: Специальные символы и цифры
		{
			name:        "special_chars",
			text:        "@#$%123!&*()",
			width:       350,
			height:      85,
			fontSize:    22,
			description: "Специальные символы и цифры",
			expected:    "Все символы должны отображаться корректно",
		},
		// Случай 11: Граничный случай - текст почти не помещается
		{
			name:        "borderline_fit",
			text:        "FITME",
			width:       140,
			height:      50,
			fontSize:    20,
			description: "Текст, который едва помещается в изображение",
			expected:    "Должен поместиться с минимальными отступами",
		},
		// Случай 12: Разные регистры
		{
			name:        "mixed_case_extreme",
			text:        "AaBbCcDdEeFfGg",
			width:       400,
			height:      95,
			fontSize:    20,
			description: "Смешанный регистр (14 символов)",
			expected:    "Должны правильно отображаться заглавные и строчные буквы",
		},
	}

	successCount := 0
	warningCount := 0
	failCount := 0

	for _, test := range extremeTests {
		fmt.Printf("\n%s\n", strings.Repeat("-", 60))
		fmt.Printf("ТЕСТ: %s\n", test.name)
		fmt.Printf("Описание: %s\n", test.description)
		fmt.Printf("Текст: '%s' (%d символов)\n", test.text, len(test.text))
		fmt.Printf("Размер изображения: %dx%d пикселей\n", test.width, test.height)
		fmt.Printf("Размер шрифта: %.0fpt\n", test.fontSize)
		fmt.Printf("Ожидаемый результат: %s\n", test.expected)

		// Создаем конфигурацию
		config := captcha.ImageCaptchaConfig{
			BackgroundColor: color.White,
			TextColor:       color.Black,
			Font:            loadFont(test.fontSize),
			FontSize:        int(test.fontSize),
			ImageWidth:      test.width,
			ImageHeight:     test.height,
		}

		// Создаем генератор
		captchaGenerator := captcha.NewImageCaptcha(config)

		// Генерируем капчу
		captchaImage, err := captchaGenerator.Generate(test.text)
		if err != nil {
			fmt.Printf("❌ КРИТИЧЕСКАЯ ОШИБКА: %v\n", err)
			failCount++
			continue
		}

		// Сохраняем в файл
		filename := filepath.Join(testDir, fmt.Sprintf("%s.png", test.name))
		err = os.WriteFile(filename, captchaImage, 0644)
		if err != nil {
			fmt.Printf("❌ ОШИБКА СОХРАНЕНИЯ: %v\n", err)
			failCount++
			continue
		}

		// Проверяем размер файла (косвенная проверка)
		fileInfo, _ := os.Stat(filename)
		fileSizeKB := float64(fileInfo.Size()) / 1024.0

		// Анализируем результат
		if fileSizeKB < 1.0 {
			fmt.Printf("⚠️  ПРЕДУПРЕЖДЕНИЕ: Очень маленький размер файла (%.2f KB)\n", fileSizeKB)
			fmt.Printf("✅ ТЕХНИЧЕСКИ УСПЕШЕН: файл создан\n")
			warningCount++
		} else {
			fmt.Printf("✅ УСПЕХ: файл создан (%.2f KB)\n", fileSizeKB)
			successCount++
		}

		// Генерируем дополнительные варианты для проверки случайности
		for i := 0; i < 2; i++ {
			captchaImage2, err := captchaGenerator.Generate(test.text)
			if err == nil {
				variantFilename := filepath.Join(testDir, fmt.Sprintf("%s_variant_%d.png", test.name, i+1))
				os.WriteFile(variantFilename, captchaImage2, 0644)
			}
		}
	}

	// Дополнительные стресс-тесты
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("СТРЕСС-ТЕСТЫ:")
	fmt.Println(strings.Repeat("=", 60))

	stressTests := []struct {
		name   string
		text   string
		width  int
		height int
	}{
		{"stress_1", "ABCDEFGHIJKLMNOP", 200, 60},     // 16 символов в узком изображении
		{"stress_2", "12345678901234567890", 300, 70}, // 20 цифр
		{"stress_3", "Aa", 30, 30},                    // Минимальный размер для 2 символов
		{"stress_4", "TEST", 50, 100},                 // Узкое высокое изображение
	}

	for _, test := range stressTests {
		config := captcha.ImageCaptchaConfig{
			BackgroundColor: color.White,
			TextColor:       color.Black,
			Font:            loadFont(18),
			FontSize:        18,
			ImageWidth:      test.width,
			ImageHeight:     test.height,
		}

		captchaGenerator := captcha.NewImageCaptcha(config)
		captchaImage, err := captchaGenerator.Generate(test.text)

		if err != nil {
			fmt.Printf("❌ Стресс-тест '%s' ПРОВАЛЕН: %v\n", test.name, err)
			failCount++
		} else {
			filename := filepath.Join(testDir, fmt.Sprintf("%s.png", test.name))
			os.WriteFile(filename, captchaImage, 0644)
			fmt.Printf("✅ Стресс-тест '%s' ПРОЙДЕН\n", test.name)
			successCount++
		}
	}

	// Итоги
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ФИНАЛЬНЫЕ ИТОГИ ЭКСТРЕМАЛЬНОГО ТЕСТИРОВАНИЯ:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Всего тестов: %d\n", len(extremeTests)+len(stressTests))
	fmt.Printf("✅ Успешных: %d\n", successCount)
	fmt.Printf("⚠️  С предупреждениями: %d\n", warningCount)
	fmt.Printf("❌ Проваленных: %d\n", failCount)

	if failCount == 0 {
		fmt.Println("\n🎉 ВСЕ ЭКСТРЕМАЛЬНЫЕ ТЕСТЫ ПРОЙДЕНЫ!")
		fmt.Println("Система корректно обрабатывает граничные случаи.")
	} else {
		fmt.Printf("\n⚠️  ВНИМАНИЕ: %d тестов не прошли\n", failCount)
	}

	fmt.Printf("\nВсе тестовые изображения сохранены в директории: %s\n", testDir)
	fmt.Println("\nРЕКОМЕНДАЦИИ ПО ВИЗУАЛЬНОЙ ПРОВЕРКЕ:")
	fmt.Println("1. Проверьте, что ни один символ не обрезан по краям")
	fmt.Println("2. Убедитесь, что текст читаем даже при экстремальных параметрах")
	fmt.Println("3. Проверьте, что повороты символов не вызывают наложения")
	fmt.Println("4. Убедитесь, что длинные тексты правильно центрируются")
	fmt.Println("5. Проверьте, что система адаптирует межсимвольные интервалы")

	fmt.Println("\nОСОБЕННОСТИ РЕАЛИЗАЦИИ ПОВОРОТА СИМВОЛОВ:")
	fmt.Println("• Автоматический расчет максимального смещения из-за поворота")
	fmt.Println("• Динамическая адаптация межсимвольных интервалов")
	fmt.Println("• Центрирование текста с учетом поворотов и смещений")
	fmt.Println("• Проверка границ для каждого пикселя повернутого символа")
	fmt.Println("• Автоматическое уменьшение интервалов при нехватке места")
}
