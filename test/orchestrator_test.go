package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/horhhe/disperser-calculator/internal/config"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/handlers"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/services"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/storage"
)

func TestOrchestratorCalculate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := storage.NewMemoryStorage()
	cfg := config.InitEnv()
	exprManager := services.NewExpressionManager(store, cfg)

	router := gin.Default()
	router.POST("/api/v1/calculate", handlers.CalculateExpression(exprManager))

	body := []byte(`{"expression": "2+2*2"}`)
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	_, exists := resp["id"]
	assert.True(t, exists, "Response should contain an id")
}
