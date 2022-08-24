package handler

import "github.com/gin-gonic/gin"

// health check
func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": map[string]string{"data": "pong"},
		"msg":  "success",
	})
}
