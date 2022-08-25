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
	"github.com/renbw/cargo/initialize"
	"github.com/spf13/viper"

	// mod file define the base package name eg `github.com/renbw/cargo`
	// `componert` is the sub package name
	log "github.com/sirupsen/logrus"
)

const (
	WAIT_TIMEOUT = time.Second * 5
)

// global conifg infomation
var conf *viper.Viper

func init() {
	// init logrus
	configLog()
	// init conifg
	conf = initialize.Conf
}

func main() {
	engine := configEngine(conf)

	v1 := engine.Group("/v1")
	{
		v1.GET("/ping", handler.Ping)
	}

	addr := conf.GetString("server.bind")
	port := conf.GetInt32("server.port")

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", addr, port), Handler: engine}
	// start in one gorotine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server start failed", err)
		}
	}()
	// wait signal to stop
	log.Info("server start at " + addr + "")

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
func configEngine(conf *viper.Viper) *gin.Engine {
	// read config file to gin engine mode
	modeKey := "app.mode"
	var mode string
	if conf.IsSet(modeKey) {
		mode = conf.GetString(modeKey)
	} else {
		mode = "release"
	}
	gin.SetMode(mode)
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
