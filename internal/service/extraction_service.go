package service

import (
	"contract-key-extractor/internal/config"
	"contract-key-extractor/internal/model"
	"contract-key-extractor/internal/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type ExtractionService struct {
	parserManager *parser.ParserManager
	aiClient      *AIServiceClient
	cfg           *config.Config
	logger        *zap.Logger
	tasks         sync.Map
}

type Task struct {
	ID          string
	Status      string
	Progress    float64
	TotalFiles  int
	Processed   int
	Failed      int
	ResultPath  string
	Error       string
	CreatedAt   time.Time
	CompletedAt time.Time
	Results     []model.ExtractionResult
}

func NewExtractionService(
	parserManager *parser.ParserManager,
	aiClient *AIServiceClient,
	cfg *config.Config,
	logger *zap.Logger,
) *ExtractionService {
	return &ExtractionService{
		parserManager: parserManager,
		aiClient:      aiClient,
		cfg:           cfg,
		logger:        logger,
	}
}

func (s *ExtractionService) ProcessFiles(filePaths []string) (*Task, error) {
	taskID := uuid.New().String()
	task := &Task{
		ID:         taskID,
		Status:     "pending",
		Progress:   0,
		TotalFiles: len(filePaths),
		CreatedAt:  time.Now(),
	}

	s.tasks.Store(taskID, task)

	go s.processTask(task, filePaths)

	return task, nil
}

func (s *ExtractionService) processTask(task *Task, filePaths []string) {
	task.Status = "processing"
	s.tasks.Store(task.ID, task)

	var results []model.ExtractionResult
	var mu sync.Mutex

	for _, filePath := range filePaths {
		result, err := s.processSingleFile(filePath)
		if err != nil {
			s.logger.Error("failed to process file",
				zap.String("file", filePath),
				zap.Error(err),
			)
			task.Failed++
		} else {
			mu.Lock()
			results = append(results, *result)
			mu.Unlock()
		}

		task.Processed++
		task.Progress = float64(task.Processed) / float64(task.TotalFiles) * 100
		s.tasks.Store(task.ID, task)
	}

	task.Results = results
	task.Status = "completed"
	task.CompletedAt = time.Now()

	outputPath, err := s.exportResults(results, task.ID)
	if err != nil {
		task.Error = fmt.Sprintf("failed to export results: %v", err)
	} else {
		task.ResultPath = outputPath
	}

	s.tasks.Store(task.ID, task)
}

func (s *ExtractionService) processSingleFile(filePath string) (*model.ExtractionResult, error) {
	startTime := time.Now()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	doc, err := s.parserManager.Parse(filePath, data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	s.logger.Info("Parsed document",
		zap.String("file", filePath),
		zap.String("fileType", string(doc.FileType)),
		zap.Bool("isScanned", doc.IsScanned),
		zap.Int("contentLen", len(doc.Content)),
	)

	var aiResp *model.AIExtractionResponse

	if doc.FileType == model.FileTypePDF {
		s.logger.Info("Calling PDF OCR", zap.String("file", filePath))
		pdfText, err := s.aiClient.PerformPDFOCR(data)
		if err != nil {
			s.logger.Warn("PDF OCR failed",
				zap.String("file", filePath),
				zap.Error(err),
			)
			doc.Content = "[PDF OCR failed, please try uploading Word or Excel format]"
		} else {
			s.logger.Info("PDF OCR succeeded", zap.Int("textLen", len(pdfText)))
			doc.Content = pdfText
			doc.IsScanned = false
		}
	} else if doc.IsScanned {
		ocrText, err := s.aiClient.PerformOCR(data)
		if err != nil {
			s.logger.Warn("OCR failed, using original content",
				zap.String("file", filePath),
				zap.Error(err),
			)
		} else {
			doc.Content = ocrText
		}
	}

	aiResp, err = s.aiClient.ExtractContractInfo(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract contract info: %w", err)
	}

	result := &model.ExtractionResult{
		ID:                uuid.New().String(),
		FileName:          filepath.Base(filePath),
		ContractInfo:      aiResp.ContractInfo,
		PartyA:            aiResp.PartyA,
		PartyB:            aiResp.PartyB,
		Financial:         aiResp.Financial,
		Validity:          aiResp.Validity,
		RightsObligations: aiResp.RightsObligations,
		BreachLiability:   aiResp.BreachLiability,
		DisputeResolution: aiResp.DisputeResolution,
		ConfidentialityIP: aiResp.ConfidentialityIP,
		OtherTerms:        aiResp.OtherTerms,
		Signature:         aiResp.Signature,
		TypeSpecific:      aiResp.TypeSpecific,
		Metadata: model.Metadata{
			SourceFile:         filePath,
			PageCount:          doc.PageCount,
			ExtractionTime:     time.Now(),
			ProcessingDuration: time.Since(startTime).Seconds(),
			OverallConfidence:  s.calculateOverallConfidence(aiResp),
			OCRRequired:        doc.IsScanned,
		},
	}

	return result, nil
}

func (s *ExtractionService) calculateOverallConfidence(resp *model.AIExtractionResponse) float64 {
	confidences := []float64{
		resp.ContractInfo.Confidence,
		resp.PartyA.Confidence,
		resp.PartyB.Confidence,
		resp.Financial.Confidence,
		resp.Validity.Confidence,
		resp.RightsObligations.Confidence,
		resp.BreachLiability.Confidence,
		resp.DisputeResolution.Confidence,
		resp.ConfidentialityIP.Confidence,
		resp.OtherTerms.Confidence,
		resp.Signature.Confidence,
	}

	var sum float64
	var count int
	for _, c := range confidences {
		if c > 0 {
			sum += c
			count++
		}
	}

	if count == 0 {
		return 0.8
	}
	return sum / float64(count)
}

func (s *ExtractionService) exportResults(results []model.ExtractionResult, taskID string) (string, error) {
	outputDir := s.cfg.Output.Path
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("extraction_result_%s.xlsx", taskID))

	return s.exportToExcel(results, outputPath)
}

func (s *ExtractionService) exportToExcel(results []model.ExtractionResult, outputPath string) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "合同提取结果"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"序号", "文件名", "合同类型", "合同编号", "签订日期", "生效日期", "到期日期",
		"甲方名称", "甲方类型", "甲方法定代表人", "甲方地址", "甲方联系方式",
		"乙方名称", "乙方类型", "乙方法定代表人", "乙方地址", "乙方联系方式",
		"交易金额", "币种", "支付方式", "付款安排",
		"生效条件", "解除条件", "合同状态",
		"甲方主要义务", "乙方主要义务", "甲方主要权利", "乙方主要权利",
		"违约情形", "违约金条款", "免责条款",
		"争议解决方式", "管辖法院", "仲裁机构", "适用法律",
		"保密条款", "知识产权归属",
		"变更条款", "转让条款", "合同份数",
		"甲方盖章", "乙方盖章",
		"置信度", "处理时间",
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", WrapText: true},
	})

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, result := range results {
		rowNum := row + 2

		contractTypeMap := map[string]string{
			"purchase":   "买卖合同",
			"lease":      "租赁合同",
			"loan":       "借款合同",
			"employment": "劳动合同",
			"service":    "服务合同",
			"other":      "其他合同",
		}
		contractType := contractTypeMap[string(result.ContractInfo.ContractType)]
		if contractType == "" {
			contractType = string(result.ContractInfo.ContractType)
		}

		data := []interface{}{
			row + 1,
			result.FileName,
			contractType,
			result.ContractInfo.ContractNumber,
			result.ContractInfo.SigningDate,
			result.ContractInfo.EffectiveDate,
			result.ContractInfo.ExpiryDate,
			result.PartyA.Name,
			result.PartyA.Type,
			result.PartyA.LegalRepresentative,
			result.PartyA.Address,
			result.PartyA.Contact,
			result.PartyB.Name,
			result.PartyB.Type,
			result.PartyB.LegalRepresentative,
			result.PartyB.Address,
			result.PartyB.Contact,
			result.Financial.TransactionAmount,
			result.Financial.Currency,
			result.Financial.PaymentMethod,
			result.Financial.PaymentSchedule,
			result.Validity.EffectiveCondition,
			result.Validity.TerminationCondition,
			result.Validity.ContractStatus,
			strings.Join(result.RightsObligations.PartyAObligations, "\n"),
			strings.Join(result.RightsObligations.PartyBObligations, "\n"),
			strings.Join(result.RightsObligations.PartyARights, "\n"),
			strings.Join(result.RightsObligations.PartyBRights, "\n"),
			strings.Join(result.BreachLiability.BreachScenarios, "\n"),
			result.BreachLiability.LiquidatedDamages,
			strings.Join(result.BreachLiability.ExemptionClauses, "\n"),
			result.DisputeResolution.ResolutionMethod,
			result.DisputeResolution.JurisdictionCourt,
			result.DisputeResolution.ArbitrationOrg,
			result.DisputeResolution.GoverningLaw,
			result.ConfidentialityIP.ConfidentialityClause,
			result.ConfidentialityIP.IPOwnership,
			result.OtherTerms.ModificationClause,
			result.OtherTerms.AssignmentClause,
			result.OtherTerms.ContractCopies,
			fmt.Sprintf("%v", result.Signature.PartyASeal),
			fmt.Sprintf("%v", result.Signature.PartyBSeal),
			fmt.Sprintf("%.1f%%", result.Metadata.OverallConfidence*100),
			fmt.Sprintf("%.2fs", result.Metadata.ProcessingDuration),
		}

		for col, value := range data {
			cell, _ := excelize.CoordinatesToCellName(col+1, rowNum)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, cellStyle)
		}
	}

	f.SetColWidth(sheetName, "A", "A", 6)
	f.SetColWidth(sheetName, "B", "B", 30)
	f.SetColWidth(sheetName, "C", "G", 15)
	f.SetColWidth(sheetName, "H", "L", 20)
	f.SetColWidth(sheetName, "M", "Q", 20)
	f.SetColWidth(sheetName, "R", "U", 15)
	f.SetColWidth(sheetName, "V", "X", 20)
	f.SetColWidth(sheetName, "Y", "AB", 30)
	f.SetColWidth(sheetName, "AC", "AE", 20)
	f.SetColWidth(sheetName, "AF", "AJ", 15)
	f.SetColWidth(sheetName, "AK", "AM", 20)
	f.SetColWidth(sheetName, "AN", "AP", 10)

	if len(results) > 0 {
		f.SetRowHeight(sheetName, 1, 30)
		for i := range results {
			f.SetRowHeight(sheetName, i+2, 60)
		}
	}

	if err := f.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("failed to save excel file: %w", err)
	}

	return outputPath, nil
}

func (s *ExtractionService) GetTaskStatus(taskID string) (*Task, error) {
	value, ok := s.tasks.Load(taskID)
	if !ok {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	return value.(*Task), nil
}

func (s *ExtractionService) GetTaskResults(taskID string) ([]model.ExtractionResult, error) {
	task, err := s.GetTaskStatus(taskID)
	if err != nil {
		return nil, err
	}
	return task.Results, nil
}
