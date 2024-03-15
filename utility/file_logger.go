package utility

import (
	"jethur/config"
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

		currentDate := time.Now().Format("02-01-2006") // DD-MM-YYYY format
		logger.SetOutput(&lumberjack.Logger{
			Filename:   "./logs/" + currentDate + ".log",
			MaxSize:    10,   // megabytes after which new file is created
			MaxBackups: 3,    // number of backups
			MaxAge:     28,   // days
			Compress:   true, // whether to compress the rotated files
			LocalTime:  true, // use local time for file naming
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
