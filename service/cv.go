package service

import (
	"fmt"
	"jobfind/cv"
)

type CVService struct {
	cvReader *cv.Reader
}

func NewCVService(cvReader *cv.Reader) *CVService {
	return &CVService{
		cvReader: cvReader,
	}
}

func (cvs *CVService) ExtractText(filePath string) (string, error) {

	result, err := cvs.cvReader.ReadPDF(filePath)

	if err != nil {
		return "", fmt.Errorf("pdf read fail: %w", err)
	}

	return result, nil
}
