package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var infoLogger *log.Logger
var debugLogger *log.Logger
var errorLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  ", log.Ltime|log.Lmicroseconds)
	debugLogger = log.New(os.Stdout, "DEBUG ", log.Ltime|log.Lmicroseconds)
	errorLogger = log.New(os.Stdout, "ERROR ", log.Ltime|log.Lmicroseconds)
}

func Info(message string) {
	appendMessage(message, infoLogger)
}

func Debug(message string) {
	appendMessage(message, debugLogger)
}

func Error(message string) {
	appendMessage(message, errorLogger)
}

func appendMessage(message string, logger *log.Logger) {
	logger.Println(message)
}

func RotateLogFiles() {
	for range time.Tick(time.Second) {
		writer := io.MultiWriter(getOutputFile(), os.Stdout)
		infoLogger.SetOutput(writer)
		debugLogger.SetOutput(writer)
		errorLogger.SetOutput(writer)
	}
}

func getOutputFile() (file *os.File) {
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			errorLogger.Println(err)
		}
	}(file)

	directory := os.Getenv("LOG_PATH")
	filename := fmt.Sprintf(
		"%s/%s.log",
		directory,
		time.Now().Format("2006-01-02"),
	)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if os.MkdirAll(directory, 0755) != nil {
			errorLogger.Println(err)
		}
		file, _ = os.Create(filename)
	} else {
		file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0755)
	}

	return
}
