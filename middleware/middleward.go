package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "服务错误",
				})
			}
		}()
		c.Next()
	}
}
