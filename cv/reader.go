package cv

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

type Reader struct {
}

func NewCVReader() *Reader {
	return &Reader{}
}

func (cvr *Reader) ReadPDF(filePath string) (string, error) {
	pdf.DebugOn = true
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

	buf.ReadFrom(b)
	content := buf.String()

	return content, nil
}
