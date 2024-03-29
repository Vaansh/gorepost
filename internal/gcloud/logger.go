package gcloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"github.com/Vaansh/gore/internal/config"
	"google.golang.org/api/option"
)

const (
	LogDirectory = "log"
)

// Log levels constants
const (
	Info      = logging.Info
	Error     = logging.Error
	Warning   = logging.Warning
	Emergency = logging.Emergency
)

var (
	client      *logging.Client
	cloudLogger *logging.Logger

	// local logging for development
	localWarningLogger *log.Logger
	localInfoLogger    *log.Logger
	localErrorLogger   *log.Logger
)

// InitLogger Initializes the logger based on configuration
func InitLogger() error {
	cfg := config.ReadLoggerConfig()

	if cfg.LocalLog {
		err := initLocalLoggers()
		if err != nil {
			return err
		}
	}

	if cfg.CloudLog {
		ctx := context.Background()
		var err error
		client, err = logging.NewClient(ctx, cfg.ProjectId, option.WithCredentialsFile(cfg.CredentialsPath))
		if err != nil {
			return err
		}
		cloudLogger = client.Logger(cfg.LogName)
	}

	return nil
}

// Helper function to create a new log file
func openLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}

// Initializes local loggers
func initLocalLoggers() error {
	currentTime := time.Now()
	run := currentTime.Format("2006-01-02|3:4:5")
	file, err := openLogFile(fmt.Sprintf("%s/%s-%s.log", LogDirectory, "info", run))
	if err != nil {
		return err
	}
	localInfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)

	file, err = openLogFile(fmt.Sprintf("%s/%s-%s.log", LogDirectory, "warning", run))
	if err != nil {
		return err
	}
	localWarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime)

	file, err = openLogFile(fmt.Sprintf("%s/%s-%s.log", LogDirectory, "error", run))
	if err != nil {
		return err
	}
	localErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)

	return nil
}

// Logs to cloud logger
func cloudLog(severity logging.Severity, format string) {
	if cloudLogger != nil {
		logEntry := logging.Entry{
			Payload:  format,
			Severity: severity,
		}
		cloudLogger.Log(logEntry)
	}
}

// LogInfo logs info message to both local and cloud loggers
func LogInfo(format string) {
	if localInfoLogger != nil {
		localInfoLogger.Println(format)
	}

	if cloudLogger != nil {
		cloudLog(Info, format)
	}
}

// LogError logs error message to both local and cloud loggers
func LogError(format string) {
	if localErrorLogger != nil {
		localErrorLogger.Println(format)
	}

	if cloudLogger != nil {
		cloudLog(Error, format)
	}
}

// LogWarning logs warning message to both local and cloud loggers
func LogWarning(format string) {
	if localWarningLogger != nil {
		localWarningLogger.Println(format)
	}

	if cloudLogger != nil {
		cloudLog(Warning, format)
	}
}

// LogFatal logs fatal message to both local and cloud loggers
func LogFatal(format string) {
	if localErrorLogger != nil {
		localErrorLogger.Println(format)
	}

	if cloudLogger != nil {
		cloudLog(Emergency, format)
	}

	log.Fatalf(format)
}
