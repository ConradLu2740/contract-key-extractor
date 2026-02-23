package parser

import (
	"bytes"
	"contract-key-extractor/internal/model"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct{}

func NewExcelParser() *ExcelParser {
	return &ExcelParser{}
}

func (p *ExcelParser) Parse(filePath string, fileData []byte) (*model.ParsedDocument, error) {
	reader, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer reader.Close()

	var contentBuilder strings.Builder
	sheets := reader.GetSheetList()
	pageCount := len(sheets)

	for sheetIdx, sheetName := range sheets {
		rows, err := reader.GetRows(sheetName)
		if err != nil {
			continue
		}

		contentBuilder.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheetName))

		for _, row := range rows {
			rowContent := strings.Join(row, "\t")
			contentBuilder.WriteString(rowContent + "\n")
		}

		if sheetIdx < len(sheets)-1 {
			contentBuilder.WriteString("\n\n")
		}
	}

	return &model.ParsedDocument{
		FileName:   filepath.Base(filePath),
		FileType:   model.FileTypeExcel,
		Content:    contentBuilder.String(),
		PageCount:  pageCount,
		IsScanned:  false,
		ImagePaths: nil,
	}, nil
}

func (p *ExcelParser) Supports(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".xlsx" || ext == ".xls"
}
