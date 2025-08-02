package utils

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(b)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func SanitizeText(input string) string {
	s := strings.ReplaceAll(input, `"`, `'`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return ' '
		}
		return r
	}, s)
	return s
}
