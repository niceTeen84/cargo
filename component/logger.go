package component

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	file, err := os.OpenFile("access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("init src error ")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log := logrus.New()
	// log info both show on std out and file appender
	writes := []io.Writer{
		file,
		os.Stdout}
	combine := io.MultiWriter(writes...)
	log.SetOutput(combine)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{})
	log.Info("logger init...")
	logger = log
}

func GetLogger() *logrus.Logger {
	return logger
}
