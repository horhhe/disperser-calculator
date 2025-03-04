package services

import (
	"errors"
	"strconv"
	"sync"

	"github.com/horhhe/disperser-calculator/internal/config"
	"github.com/horhhe/disperser-calculator/internal/models"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/storage"
)

type ExpressionManager struct {
	store         storage.Storage
	cfg           *config.EnvConfig
	mu            sync.Mutex
	taskIDCounter int
}

func NewExpressionManager(store storage.Storage, cfg *config.EnvConfig) *ExpressionManager {
	return &ExpressionManager{
		store: store,
		cfg:   cfg,
	}
}

func (em *ExpressionManager) CreateExpression(expr string) (string, error) {
	if expr == "" {
		return "", errors.New("empty expression")
	}

	expressionID := em.store.CreateExpression(expr)
	em.mu.Lock()
	em.taskIDCounter++
	taskID := em.taskIDCounter
	em.mu.Unlock()

	task := models.Task{
		ID:            taskID,
		ExpressionID:  expressionID,
		Arg1:          expr,
		Operation:     "eval",
		OperationTime: em.cfg.TimeEvaluation,
		Status:        "pending",
	}

	em.store.AddTask(expressionID, task)
	return expressionID, nil
}

func (em *ExpressionManager) GetAllExpressions() []models.Expression {
	return em.store.GetExpressions()
}

func (em *ExpressionManager) GetExpressionByID(id string) (models.Expression, bool) {
	return em.store.GetExpression(id)
}

func (em *ExpressionManager) GetNextTask() (models.Task, bool) {
	return em.store.GetPendingTask()
}

func (em *ExpressionManager) SetTaskResult(taskID string, result float64) error {
	tid, _ := strconv.Atoi(taskID)
	return em.store.CompleteTask(tid, result)
}
