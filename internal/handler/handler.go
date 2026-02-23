package handler

import (
	"contract-key-extractor/internal/model"
	"contract-key-extractor/internal/service"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	extractionService *service.ExtractionService
	uploadPath        string
	logger            *zap.Logger
}

func NewHandler(extractionService *service.ExtractionService, uploadPath string, logger *zap.Logger) *Handler {
	return &Handler{
		extractionService: extractionService,
		uploadPath:        uploadPath,
		logger:            logger,
	}
}

func (h *Handler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	if err := os.MkdirAll(h.uploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	var filePaths []string
	for _, file := range files {
		dst := filepath.Join(h.uploadPath, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			h.logger.Error("failed to save file", zap.String("file", file.Filename), zap.Error(err))
			continue
		}
		filePaths = append(filePaths, dst)
	}

	task, err := h.extractionService.ProcessFiles(filePaths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id":     task.ID,
		"status":      task.Status,
		"total_files": task.TotalFiles,
		"message":     "files uploaded successfully, processing started",
	})
}

func (h *Handler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")

	task, err := h.extractionService.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := model.TaskStatus{
		TaskID:     task.ID,
		Status:     task.Status,
		Progress:   task.Progress,
		TotalFiles: task.TotalFiles,
		Processed:  task.Processed,
		Failed:     task.Failed,
		ResultPath: task.ResultPath,
		Error:      task.Error,
		CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if !task.CompletedAt.IsZero() {
		response.CompletedAt = task.CompletedAt.Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetTaskResults(c *gin.Context) {
	taskID := c.Param("task_id")

	results, err := h.extractionService.GetTaskResults(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ExtractionResponse{
		Success: true,
		Message: "extraction completed",
		Results: results,
	})
}

func (h *Handler) DownloadResult(c *gin.Context) {
	taskID := c.Param("task_id")

	task, err := h.extractionService.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if task.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not completed yet"})
		return
	}

	if task.ResultPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "result file not found"})
		return
	}

	c.FileAttachment(task.ResultPath, filepath.Base(task.ResultPath))
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "contract-key-extractor",
	})
}
