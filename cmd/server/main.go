package main

import (
	"contract-key-extractor/internal/config"
	"contract-key-extractor/internal/handler"
	"contract-key-extractor/internal/parser"
	"contract-key-extractor/internal/service"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := config.InitLogger(&cfg.Logging)
	if err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}

	parserManager := parser.NewParserManager(logger)

	aiClient := service.NewAIServiceClient(&cfg.AIService, logger)

	extractionService := service.NewExtractionService(parserManager, aiClient, cfg, logger)

	h := handler.NewHandler(extractionService, cfg.Upload.Path, logger)

	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1")
	{
		api.POST("/upload", h.UploadFiles)
		api.GET("/task/:task_id", h.GetTaskStatus)
		api.GET("/task/:task_id/results", h.GetTaskResults)
		api.GET("/task/:task_id/download", h.DownloadResult)
	}

	router.GET("/health", h.HealthCheck)

	router.Static("/static", "./web/dist")
	router.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("starting server", zap.String("address", addr))

	if err := router.Run(addr); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
