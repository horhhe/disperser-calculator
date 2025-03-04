package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/horhhe/disperser-calculator/internal/models"
)

type memoryStorage struct {
	mu           sync.Mutex
	expressions  map[string]models.Expression
	tasks        map[int]models.Task
	expressionID int
}

type Storage interface {
	CreateExpression(expr string) string
	GetExpression(id string) (models.Expression, bool)
	GetExpressions() []models.Expression

	AddTask(exprID string, task models.Task)
	GetPendingTask() (models.Task, bool)
	CompleteTask(taskID int, result float64) error
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		expressions: make(map[string]models.Expression),
		tasks:       make(map[int]models.Task),
	}
}

func (m *memoryStorage) CreateExpression(expr string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expressionID++
	id := fmt.Sprintf("%d", m.expressionID)
	m.expressions[id] = models.Expression{
		ID:     id,
		Status: "pending",
		Result: 0,
	}
	return id
}

func (m *memoryStorage) GetExpression(id string) (models.Expression, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	exp, ok := m.expressions[id]
	return exp, ok
}

func (m *memoryStorage) GetExpressions() []models.Expression {
	m.mu.Lock()
	defer m.mu.Unlock()

	res := make([]models.Expression, 0, len(m.expressions))
	for _, v := range m.expressions {
		res = append(res, v)
	}
	return res
}

func (m *memoryStorage) AddTask(exprID string, task models.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[task.ID] = task
}

func (m *memoryStorage) GetPendingTask() (models.Task, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, task := range m.tasks {
		if task.Status == "pending" {
			// Обновим статус, чтобы его не взял другой агент
			task.Status = "in-progress"
			m.tasks[id] = task
			return task, true
		}
	}
	return models.Task{}, false
}

func (m *memoryStorage) CompleteTask(taskID int, result float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[taskID]
	if !ok {
		return errors.New("task not found")
	}

	task.Status = "done"
	m.tasks[taskID] = task

	expr, ok := m.expressions[task.ExpressionID]
	if !ok {
		return errors.New("expression not found")
	}
	expr.Status = "done"
	expr.Result = result
	m.expressions[task.ExpressionID] = expr
	return nil
}
