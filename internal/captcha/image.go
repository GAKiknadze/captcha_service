package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageCaptchaConfig struct {
	BackgroundColor color.Color
	TextColor       color.Color
	Font            *font.Face
	FontSize        int
	ImageWidth      int
	ImageHeight     int
}

type ImageCaptcha struct {
	backgroundColor color.Color
	textColor       color.Color
	font            *font.Face
	fontSize        int
	imageWidth      int
	imageHeight     int
}

func NewImageCaptcha(config ImageCaptchaConfig) *ImageCaptcha {
	return &ImageCaptcha{
		backgroundColor: config.BackgroundColor,
		textColor:       config.TextColor,
		font:            config.Font,
		fontSize:        config.FontSize,
		imageWidth:      config.ImageWidth,
		imageHeight:     config.ImageHeight,
	}
}

func (c *ImageCaptcha) Generate(code string) ([]byte, error) {
	// Создаем изображение с заданными размерами
	width := c.imageWidth
	height := c.imageHeight
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Заполняем фон
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, c.backgroundColor)
		}
	}

	// Рисуем текст капчи с поворотом и смещением символов
	rand.Seed(time.Now().UnixNano())

	// Рассчитываем параметры для правильного позиционирования
	charWidth := 35
	charHeight := 45
	maxRotation := 20.0 // Максимальный угол поворота в градусах

	// Максимальное смещение из-за поворота (диагональ символа * sin(угла))
	maxRotationRad := maxRotation * math.Pi / 180
	diagonal := math.Sqrt(float64(charWidth*charWidth + charHeight*charHeight))
	maxRotationOffset := int(diagonal * math.Sin(maxRotationRad))

	// Рассчитываем общую ширину текста с учетом поворотов и случайных интервалов
	// Базовый интервал между символами
	baseCharSpacing := 25
	// Максимальная дополнительная вариация интервала
	maxSpacingVariation := 5

	// Рассчитываем максимальную возможную ширину
	maxTotalWidth := len(code)*charWidth + (len(code)-1)*(baseCharSpacing+maxSpacingVariation)

	// Берем максимальную ширину для безопасного расчета
	totalTextWidth := maxTotalWidth
	totalWidthWithRotation := totalTextWidth + 2*maxRotationOffset

	// Центрируем текст по горизонтали с проверкой границ
	startX := (width - totalWidthWithRotation) / 2
	if startX < maxRotationOffset+10 {
		startX = maxRotationOffset + 10 // Минимальный отступ с учетом поворота
	}

	// Проверяем, не выходит ли текст за правую границу
	if startX+totalWidthWithRotation > width-10 {
		// Если выходит, уменьшаем начальную позицию
		startX = width - totalWidthWithRotation - 10
		if startX < maxRotationOffset+10 {
			// Если все равно не помещается, уменьшаем межсимвольные интервалы
			baseCharSpacing = 20
			maxSpacingVariation = 3
			// Пересчитываем
			maxTotalWidth = len(code)*charWidth + (len(code)-1)*(baseCharSpacing+maxSpacingVariation)
			totalTextWidth = maxTotalWidth
			totalWidthWithRotation = totalTextWidth + 2*maxRotationOffset
			startX = (width - totalWidthWithRotation) / 2
			if startX < maxRotationOffset+10 {
				startX = maxRotationOffset + 10
			}
		}
	}

	charSpacing := baseCharSpacing
	centerY := height / 2

	// Проверяем вертикальные границы
	minY := charHeight/2 + maxRotationOffset
	maxY := height - charHeight/2 - maxRotationOffset
	if centerY < minY {
		centerY = minY
	} else if centerY > maxY {
		centerY = maxY
	}

	for i, ch := range code {
		// Создаем временное изображение для символа
		charImg := image.NewRGBA(image.Rect(0, 0, charWidth, charHeight))

		// Заполняем прозрачным фоном
		for x := 0; x < charWidth; x++ {
			for y := 0; y < charHeight; y++ {
				charImg.Set(x, y, color.Transparent)
			}
		}

		// Рисуем символ во временном изображении
		charDrawer := &font.Drawer{
			Dst:  charImg,
			Src:  image.NewUniform(c.textColor),
			Face: *c.font,
			Dot:  fixed.P(8, charHeight/2+int(c.fontSize)/2),
		}
		charDrawer.DrawString(string(ch))

		// Применяем случайный поворот (-20 до +20 градусов)
		angle := (rand.Float64()*40 - 20) * math.Pi / 180 // -20° to +20° in radians

		// Добавляем случайное вертикальное смещение (-5 до +5 пикселей)
		verticalOffset := rand.Intn(11) - 5 // -5 to +5

		// Вычисляем позицию для вставки повернутого символа
		// Учитываем дополнительное пространство для поворота
		posX := startX + maxRotationOffset + i*charSpacing
		posY := centerY + verticalOffset

		// Вставляем повернутый символ в основное изображение
		for x := 0; x < charWidth; x++ {
			for y := 0; y < charHeight; y++ {
				// Получаем цвет пикселя из символа
				r, g, b, a := charImg.At(x, y).RGBA()
				if a > 0 {
					// Вычисляем координаты относительно центра символа
					relX := float64(x - charWidth/2)
					relY := float64(y - charHeight/2)

					// Применяем вращение
					rotX := relX*math.Cos(angle) - relY*math.Sin(angle)
					rotY := relX*math.Sin(angle) + relY*math.Cos(angle)

					// Возвращаем к абсолютным координатам с учетом смещения
					destX := int(rotX) + posX + charWidth/2
					destY := int(rotY) + posY

					// Проверяем границы и рисуем пиксель
					if destX >= 0 && destX < width && destY >= 0 && destY < height {
						// Также проверяем, что пиксель не слишком близко к краю
						// (оставляем запас в 2 пикселя для искажений)
						if destX >= 2 && destX < width-2 && destY >= 2 && destY < height-2 {
							img.Set(destX, destY, color.RGBA{
								R: uint8(r >> 8),
								G: uint8(g >> 8),
								B: uint8(b >> 8),
								A: uint8(a >> 8),
							})
						}
					}
				}
			}
		}

		// Добавляем небольшую случайную вариацию в межсимвольный интервал
		spacingVariation := rand.Intn(2*maxSpacingVariation+1) - maxSpacingVariation // -maxSpacingVariation to +maxSpacingVariation
		charSpacing = baseCharSpacing + spacingVariation

		// Проверяем, не выйдет ли следующий символ за границы
		if i < len(code)-1 {
			nextPosX := posX + charWidth + charSpacing
			if nextPosX+charWidth+maxRotationOffset > width-10 {
				// Уменьшаем интервал для последующих символов
				charSpacing = 15
				maxSpacingVariation = 2
			}
		}
	}

	// Добавляем случайные помехи - точки
	for i := 0; i < 100; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		})
	}

	// Добавляем случайные линии
	for i := 0; i < 5; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)

		lineColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: uint8(rand.Intn(100) + 100),
		}

		// Простая реализация линии
		dx := abs(x2 - x1)
		dy := abs(y2 - y1)
		sx := -1
		if x1 < x2 {
			sx = 1
		}
		sy := -1
		if y1 < y2 {
			sy = 1
		}
		err := dx - dy

		for {
			if x1 >= 0 && x1 < width && y1 >= 0 && y1 < height {
				img.Set(x1, y1, lineColor)
			}
			if x1 == x2 && y1 == y2 {
				break
			}
			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x1 += sx
			}
			if e2 < dx {
				err += dx
				y1 += sy
			}
		}
	}

	// Добавляем волнообразное искажение текста
	distorted := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		// Более сложное искажение с двумя синусоидами
		offset1 := int(3 * math.Sin(float64(x)*0.08))
		offset2 := int(2 * math.Sin(float64(x)*0.15+1.5))
		offset := offset1 + offset2

		for y := 0; y < height; y++ {
			srcY := y + offset
			if srcY >= 0 && srcY < height {
				distorted.Set(x, y, img.At(x, srcY))
			} else {
				distorted.Set(x, y, c.backgroundColor)
			}
		}
	}

	// Добавляем легкое горизонтальное искажение
	finalImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		// Небольшое горизонтальное смещение
		horizOffset := int(2 * math.Sin(float64(y)*0.1))
		for x := 0; x < width; x++ {
			srcX := x + horizOffset
			if srcX >= 0 && srcX < width {
				finalImage.Set(x, y, distorted.At(srcX, y))
			} else {
				finalImage.Set(x, y, c.backgroundColor)
			}
		}
	}

	// Кодируем изображение в PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, finalImage)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Вспомогательная функция для абсолютного значения
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
