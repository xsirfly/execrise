package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortRequesrtWithError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}
