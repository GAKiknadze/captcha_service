package captcha

type Captcha interface {
	Generate(code string) ([]byte, error)
}
