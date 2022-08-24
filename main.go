package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/renbw/cargo/handler"

	// mod file define the base package name eg `github.com/renbw/cargo`
	// `componert` is the sub package name
	log "github.com/sirupsen/logrus"
)

const (
	BIND_ADDR    = "0.0.0.0"
	PORT         = 8080
	WAIT_TIMEOUT = time.Second * 5
)

func main() {
	configLog()
	gin.SetMode(gin.DebugMode)
	engine := configEngine()

	v1 := engine.Group("/v1")
	{
		v1.GET("/ping", handler.Ping)
	}

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", BIND_ADDR, PORT), Handler: engine}
	// start in one gorotine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server start failed", err)
		}
	}()
	// wait signal to stop
	log.Info("server start at " + BIND_ADDR + "")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// shutdown gracefully

	ctx, cancel := context.WithTimeout(context.Background(), WAIT_TIMEOUT)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("show down server failed", err)
	}
	<-ctx.Done()
	log.Info("time out")
	log.Info("server exit")
	os.Exit(0)
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
