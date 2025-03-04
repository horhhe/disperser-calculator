package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horhhe/disperser-calculator/internal/orchestrator/services"
)

func GetExpressions(exprManager *services.ExpressionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		expressions := exprManager.GetAllExpressions()
		c.JSON(http.StatusOK, gin.H{"expressions": expressions})
	}
}

func GetExpressionByID(exprManager *services.ExpressionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		expr, found := exprManager.GetExpressionByID(id)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "expression not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"expression": expr})
	}
}
