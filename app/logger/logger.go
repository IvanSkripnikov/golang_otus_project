package logger

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

const serviceNameDefault = "banner-rotation"

var errorLevels = map[logrus.Level]string{
	logrus.PanicLevel: "panic",
	logrus.FatalLevel: "fatal",
	logrus.ErrorLevel: "error",
	logrus.WarnLevel:  "warning",
	logrus.InfoLevel:  "info",
	logrus.DebugLevel: "debug",
	logrus.TraceLevel: "trace",
}

func SendToPanicLog(message string) {
	pushLogger(message, logrus.PanicLevel)
	os.Exit(1)
}

func SendToFatalLog(message string) {
	pushLogger(message, logrus.FatalLevel)
	os.Exit(1)
}

func SendToErrorLog(message string) {
	pushLogger(message, logrus.ErrorLevel)
}

func SendToWarningLog(message string) {
	pushLogger(message, logrus.WarnLevel)
}

func SendToInfoLog(message string) {
	pushLogger(message, logrus.InfoLevel)
}

func SendToDebugLog(message string) {
	pushLogger(message, logrus.DebugLevel)
}

func SendToTraceLog(message string) {
	pushLogger(message, logrus.TraceLevel)
}

func pushLogger(message string, currentLevel logrus.Level) {
	configLogLevel := os.Getenv("LOG_LEVEL")

	if len(configLogLevel) == 0 {
		configLogLevel = "2"
	}

	levelValue, errLevel := strconv.Atoi(configLogLevel)
	var logLevel logrus.Level

	if errLevel != nil {
		log.Println(errLevel)
	} else {
		logLevel = logrus.Level(levelValue)
	}

	if currentLevel > logLevel {
		return
	}

	flag.Parse()
	logsFilePath := getLogFilePath()
	logFile, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, logFile),
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			LogFormat:       "[%time%] %msg%",
		},
	}

	levelMessage := errorLevels[currentLevel]
	logger.Printf("[%s] [%s] [%s] %s \n",
		getHostName(), serviceNameDefault, levelMessage, message)
}

func getLogFilePath() string {
	containerName := os.Getenv("CONTAINER_NAME")

	if len(containerName) == 0 {
		containerName = serviceNameDefault
	}

	return fmt.Sprintf("./log/%s.log", containerName)
}

func getHostName() string {
	var hostName string
	hostNameFile, err := os.ReadFile("/etc/hostname")
	if err != nil {
		serverName, _ := os.Hostname()
		hostName = serverName
	} else {
		hostName = strings.ReplaceAll(string(hostNameFile), "\n", "")
	}

	return hostName
}
