package logger

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	serviceNameDefault = "banner-rotation"

	panicLevel = 0

	fatalLevel = 1

	errorLevel = 2

	warnLevel = 3

	infoLevel = 4

	debugLevel = 5

	traceLevel = 6
)

var errorLevels = map[int]string{
	panicLevel: "panic",

	fatalLevel: "fatal",

	errorLevel: "error",

	warnLevel: "warning",

	infoLevel: "info",

	debugLevel: "debug",

	traceLevel: "trace",
}

func SendToPanicLog(message string) {
	pushLogger(message, panicLevel)

	os.Exit(1)
}

func SendToFatalLog(message string) {
	pushLogger(message, fatalLevel)

	os.Exit(1)
}

func SendToErrorLog(message string) {
	pushLogger(message, errorLevel)
}

func SendToWarningLog(message string) {
	pushLogger(message, warnLevel)
}

func SendToInfoLog(message string) {
	pushLogger(message, infoLevel)
}

func SendToDebugLog(message string) {
	pushLogger(message, debugLevel)
}

func SendToTraceLog(message string) {
	pushLogger(message, traceLevel)
}

func pushLogger(message string, currentLevel int) {
	configLogLevel := os.Getenv("LOG_LEVEL")

	if len(configLogLevel) == 0 {
		configLogLevel = "2"
	}

	levelValue, errLevel := strconv.Atoi(configLogLevel)

	var logLevel int

	if errLevel != nil {
		log.Println(errLevel)
	} else {
		logLevel = levelValue
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

	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	levelMessage := errorLevels[currentLevel]

	fmt.Printf("[%s] [%s] [%s] %s \n", getHostName(), serviceNameDefault, levelMessage, message)
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
