package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/horhhe/disperser-calculator/internal/config"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/handlers"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/services"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/storage"
)

func main() {
	cfg := config.InitEnv()
	store := storage.NewMemoryStorage()
	exprManager := services.NewExpressionManager(store, cfg)
	r := gin.Default()
	public := r.Group("/api/v1")
	{
		public.POST("/calculate", handlers.CalculateExpression(exprManager))
		public.GET("/expressions", handlers.GetExpressions(exprManager))
		public.GET("/expressions/:id", handlers.GetExpressionByID(exprManager))
	}

	internal := r.Group("/internal")
	{
		internal.GET("/task", handlers.GetTask(exprManager))
		internal.POST("/task", handlers.PostTaskResult(exprManager))
	}

	addr := ":8080"
	if val := os.Getenv("ORCHESTRATOR_PORT"); val != "" {
		addr = ":" + val
	}
	log.Printf("Orchestrator listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run orchestrator: %v", err)
	}
}
