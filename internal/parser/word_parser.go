package parser

import (
	"archive/zip"
	"bytes"
	"contract-key-extractor/internal/model"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type WordParser struct{}

func NewWordParser() *WordParser {
	return &WordParser{}
}

func (p *WordParser) Parse(filePath string, fileData []byte) (*model.ParsedDocument, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".docx" {
		return p.parseDocx(filePath, fileData)
	}

	return nil, fmt.Errorf("unsupported word format: %s", ext)
}

func (p *WordParser) parseDocx(filePath string, fileData []byte) (*model.ParsedDocument, error) {
	reader, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open docx file: %w", err)
	}

	var contentBuilder strings.Builder
	var documentXML []byte

	for _, file := range reader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open document.xml: %w", err)
			}
			documentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read document.xml: %w", err)
			}
			break
		}
	}

	if documentXML == nil {
		return nil, fmt.Errorf("document.xml not found in docx")
	}

	text := p.extractTextFromXML(documentXML)
	contentBuilder.WriteString(text)

	return &model.ParsedDocument{
		FileName:   filepath.Base(filePath),
		FileType:   model.FileTypeWord,
		Content:    contentBuilder.String(),
		PageCount:  1,
		IsScanned:  false,
		ImagePaths: nil,
	}, nil
}

func (p *WordParser) extractTextFromXML(data []byte) string {
	var result strings.Builder

	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				result.WriteString(text)
			}
		case xml.EndElement:
			if t.Name.Local == "p" {
				result.WriteString("\n")
			} else if t.Name.Local == "t" {
				result.WriteString(" ")
			}
		}
	}

	return result.String()
}

func (p *WordParser) Supports(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".docx" || ext == ".doc"
}
