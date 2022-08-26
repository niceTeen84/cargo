package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// health check
func Ping(ctx *gin.Context) {
	name := ctx.Query("name")
	log.Info("name is ", name)
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": map[string]string{"data": "pong"},
		"msg":  "success",
	})
}
