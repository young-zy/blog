package controllers

import "github.com/gin-gonic/gin"

func generateTraceId(c *gin.Context) {
	//c.Set("traceId", "test")
	c.Next()
}
