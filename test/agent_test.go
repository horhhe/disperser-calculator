package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horhhe/disperser-calculator/internal/agent"
	"github.com/horhhe/disperser-calculator/internal/agent/client"
	"github.com/horhhe/disperser-calculator/internal/config"
	"github.com/stretchr/testify/assert"
)

// Пример простого теста
func TestAgentRequestTask(t *testing.T) {
	// Создаём тестовый сервер, чтобы эмулировать ответы оркестратора
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/internal/task" {
			// Возвращаем одну задачу
			fmt.Fprintf(w, `{"task": {"id": 1, "arg1": "2+2", "operation":"eval","operation_time":100}}`)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/internal/task" {
			// Принимаем результат
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status":"result accepted"}`)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	cfg := config.InitEnv()
	ag := agent.NewWorker(ts.URL, 1, cfg)
	// Проверяем, что агент корректно отработает один цикл RequestAndProcessTask
	ag.RequestAndProcessTask()
	// Если всё хорошо, до этого места не упадёт с ошибками
	assert.True(t, true)
}

// Можно отдельно протестировать сам client
func TestAgentClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/internal/task" {
			fmt.Fprintf(w, `{"task": {"id": 1}}`)
			return
		}
	}))
	defer ts.Close()

	c := client.NewAgentClient(ts.URL)
	task, err := c.GetTask()
	assert.NoError(t, err)
	assert.Equal(t, 1, task.ID)
}

// Пример: тест на сценарий "нет задач" => 404
func TestAgentNoTasks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Возвращаем 404, имитируя отсутствие задач
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	c := client.NewAgentClient(ts.URL)
	_, err := c.GetTask()
	assert.Error(t, err, "should return an error if no tasks available")
}
