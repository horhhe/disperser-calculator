package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horhhe/disperser-calculator/internal/models"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/services"
)

func GetTask(exprManager *services.ExpressionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				c.JSON(http.StatusNotFound, gin.H{"error": "no tasks available"})
				return
			case <-ticker.C:
				task, ok := exprManager.GetNextTask()
				if ok {
					c.JSON(http.StatusOK, gin.H{"task": task})
					return
				}
			}
		}
	}
}

func PostTaskResult(exprManager *services.ExpressionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result models.TaskResultRequest
		if err := c.ShouldBindJSON(&result); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request"})
			return
		}
		err := exprManager.SetTaskResult(strconv.Itoa(result.ID), result.Result)
		if err != nil {
			if err.Error() == "task not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "result accepted"})
	}
}
