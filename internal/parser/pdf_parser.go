package parser

import (
	"contract-key-extractor/internal/model"
	"path/filepath"
	"strings"
)

type PDFParser struct{}

func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

func (p *PDFParser) Parse(filePath string, fileData []byte) (*model.ParsedDocument, error) {
	return &model.ParsedDocument{
		FileName:   filepath.Base(filePath),
		FileType:   model.FileTypePDF,
		Content:    "",
		PageCount:  1,
		IsScanned:  true,
		ImagePaths: nil,
	}, nil
}

func (p *PDFParser) Supports(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".pdf"
}
