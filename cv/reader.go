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

	return content, nil
}
