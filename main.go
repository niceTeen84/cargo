package main

import (
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// mod file define the base package name eg `github.com/renbw/cargo`
	// `componert` is the sub package name
	log "github.com/sirupsen/logrus"
)

func main() {
	configLog()
	gin.SetMode(gin.DebugMode)
	engine := configEngine()

	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	engine.Run(":8080")
}

// logrus config
func LogToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		// End time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}

}

// config the gin framework instance
// includ midware cors time statistics and 404 cover
func configEngine() *gin.Engine {
	r := gin.New()
	corsDefaultConf := cors.DefaultConfig()
	corsDefaultConf.AllowAllOrigins = true
	corsDefaultConf.AddAllowHeaders("Authorization")

	r.Use(cors.New(corsDefaultConf), gin.Recovery(), LogToFile())
	return r
}

// conifg thr logrus global instance
// include appender file and stdout
func configLog() {
	file, err := os.OpenFile("access.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatal("init file log failed", err.Error())
	}
	combine := io.MultiWriter(os.Stdout, file)
	log.SetOutput(combine)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
