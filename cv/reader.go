package cv

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Reader struct {
}

func NewCVReader() *Reader {
	return &Reader{}
}

func standardizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\n", "")

	return strings.Join(strings.Fields(s), " ")
}

func normalizeText(s string) string {
	s = strings.ToLower(s)

	return s
}

func cleanupText(s string) string {
	s = strings.ReplaceAll(s, "|", " ")
	s = strings.ReplaceAll(s, "•", " ")
	s = standardizeWhitespace(s)

	return s
}

func (cvr *Reader) ReadPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)

	if err != nil {
		return "", fmt.Errorf("pdf file open failed: %w", err)
	}

	defer f.Close()

	var buf bytes.Buffer

	b, err := r.GetPlainText()

	if err != nil {
		return "", fmt.Errorf("pdf plain text translation fail: %w", err)
	}

	if _, err := buf.ReadFrom(b); err != nil {
		return "", fmt.Errorf("read extracted pdf text: %w", err)
	}

	content := buf.String()

	content = cleanupText(content)
	content = normalizeText(content)

	return content, nil
}
