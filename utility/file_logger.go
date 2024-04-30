package utility

import (
	"jethur/config"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

func getLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})

		storagePath := config.GetEnv("STORAGE_PATH", "./storage")
		logsPath := filepath.Join(storagePath, "logs")
		currentDate := time.Now().Format("02-01-2006")
		logFilename := filepath.Join(logsPath, currentDate+".log")

		logger.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
			LocalTime:  true,
		})
	})

	// Set log level based on environment
	if config.GetEnv("ENV", "") == "Development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
func LogDebug(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Debug(message)
}

func LogInfo(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Info(message)
}

func LogWarning(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Warn(message)
}

func LogError(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Error(message)
}

func LogFatal(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Fatal(message)
}

func Log(message string, fields map[string]interface{}) {
	log := getLogger()
	log.WithFields(fields).Error(message)
}
