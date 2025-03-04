package agent

import (
	"log"
	"time"

	"github.com/horhhe/disperser-calculator/internal/agent/client"
	"github.com/horhhe/disperser-calculator/internal/config"
	"github.com/horhhe/disperser-calculator/internal/models"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/services"
)

type Worker struct {
	orchestratorURL string
	computingPower  int
	cfg             *config.EnvConfig
	client          *client.AgentClient
}

func NewWorker(url string, cp int, cfg *config.EnvConfig) *Worker {
	return &Worker{
		orchestratorURL: url,
		computingPower:  cp,
		cfg:             cfg,
		client:          client.NewAgentClient(url),
	}
}

func (w *Worker) RequestAndProcessTask() {
	task, err := w.client.GetTask()
	if err != nil {
		log.Println("No tasks or error:", err)
		return
	}

	log.Printf("Received task ID=%d, expr=%s\n", task.ID, task.Arg1)
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	result, err := services.EvaluateExpression(task.Arg1)
	if err != nil {
		log.Printf("Error evaluating expression: %v\n", err)
		return
	}
	if err := w.client.PostTaskResult(models.TaskResultRequest{
		ID:     task.ID,
		Result: result,
	}); err != nil {
		log.Printf("Error posting task result: %v\n", err)
		return
	}
	log.Printf("Task ID=%d completed, result=%.2f\n", task.ID, result)
}
