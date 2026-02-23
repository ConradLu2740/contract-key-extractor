package parser

import (
	"contract-key-extractor/internal/model"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type DocumentParser interface {
	Parse(filePath string, fileData []byte) (*model.ParsedDocument, error)
	Supports(filePath string) bool
}

type ParserManager struct {
	parsers []DocumentParser
	logger  *zap.Logger
}

func NewParserManager(logger *zap.Logger) *ParserManager {
	return &ParserManager{
		parsers: []DocumentParser{
			NewPDFParser(),
			NewExcelParser(),
			NewWordParser(),
		},
		logger: logger,
	}
}

func (m *ParserManager) Parse(filePath string, fileData []byte) (*model.ParsedDocument, error) {
	for _, parser := range m.parsers {
		if parser.Supports(filePath) {
			m.logger.Debug("using parser for file",
				zap.String("file", filePath),
				zap.String("parser", fmt.Sprintf("%T", parser)),
			)
			return parser.Parse(filePath, fileData)
		}
	}
	return nil, errors.New("no suitable parser found for file: " + filePath)
}

func (m *ParserManager) ParseBatch(files map[string][]byte) ([]*model.ParsedDocument, error) {
	var (
		results []*model.ParsedDocument
		mu      sync.Mutex
		wg      sync.WaitGroup
		errs    []error
	)

	for filePath, data := range files {
		wg.Add(1)
		go func(fp string, d []byte) {
			defer wg.Done()

			doc, err := m.Parse(fp, d)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("failed to parse %s: %w", fp, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			results = append(results, doc)
			mu.Unlock()
		}(filePath, data)
	}

	wg.Wait()

	if len(errs) > 0 {
		return results, fmt.Errorf("batch parsing completed with %d errors", len(errs))
	}

	return results, nil
}

func (m *ParserManager) GetSupportedExtensions() []string {
	return []string{".pdf", ".xlsx", ".xls", ".docx", ".doc"}
}

func (m *ParserManager) IsSupported(filePath string) bool {
	for _, parser := range m.parsers {
		if parser.Supports(filePath) {
			return true
		}
	}
	return false
}
