package log

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/env"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var infoLogger *log.Logger
var debugLogger *log.Logger
var errorLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  ", log.Ltime|log.Ltime)
	debugLogger = log.New(os.Stdout, "DEBUG ", log.Ltime|log.Ltime)
	errorLogger = log.New(os.Stdout, "ERROR ", log.Ltime|log.Ltime)
}

func Info(a ...any) {
	infoLogger.Println(a...)
}

func Debug(value any) {
	debugLogger.Println(value)
}

func Error(error error) {
	_, file, line, _ := runtime.Caller(1)
	relativePath, _ := filepath.Rel(env.Get("APP_PATH"), file)
	location := fmt.Sprintf("%s:%d", relativePath, line)
	errorLogger.Println(location, error)
}

func RotateFiles() {
	setOutputFile()
	for range time.Tick(time.Second) {
		if time.Now().Hour() == 0 && time.Now().Minute() == 0 {
			setOutputFile()
		}
	}
}

func setOutputFile() {
	writer := io.MultiWriter(getOutputFile(), os.Stdout)
	infoLogger.SetOutput(writer)
	debugLogger.SetOutput(writer)
	errorLogger.SetOutput(writer)
}

func getOutputFile() (file *os.File) {
	defer file.Close()

	directory := os.Getenv("LOG_PATH")
	filename := fmt.Sprintf(
		"%s/%s.log",
		directory,
		time.Now().Format("2006-01-02"),
	)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if os.MkdirAll(directory, 0666) != nil {
			errorLogger.Println(err)
		}
		file, _ = os.Create(filename)
	} else {
		file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	}

	return
}
