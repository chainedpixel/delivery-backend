package logs

import (
	"config"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	errPackage "shared/error"
	"strings"
)

type CustomFormatter struct {
	logrus.TextFormatter
}

var Logger *logrus.Logger

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColors := map[logrus.Level]string{
		logrus.DebugLevel: "\033[37m", // Tonalidad gris para debug
		logrus.InfoLevel:  "\033[32m", // Tonalidad verde para info
		logrus.WarnLevel:  "\033[33m", // Tonalidad amarilla para warning
		logrus.ErrorLevel: "\033[31m", // Tonalidad roja para error
		logrus.FatalLevel: "\033[35m", // Tonalidad morada para fatal
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	color := levelColors[entry.Level]

	baseMessage := fmt.Sprintf("%s [%s%s\033[0m] %s",
		color,
		timestamp,
		level,
		entry.Message,
	)

	var fields string
	if len(entry.Data) > 0 {
		fields = "\nDetails:\n"
		for k, v := range entry.Data {
			fields += fmt.Sprintf("	%-20s: %v\n", k, v)
		}
	}

	return []byte(fmt.Sprintf("%s%s\n", baseMessage, fields)), nil
}

func InitLogger(envConfig *config.EnvConfig) error {
	_, b, _, _ := runtime.Caller(0)
	Logger = logrus.New()

	Logger.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	})

	Logger.SetLevel(determineLogLevel(envConfig.Log.Level))
	projectRoot := filepath.Join(filepath.Dir(b))
	logFile, err := os.OpenFile(projectRoot+"system.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return errPackage.ErrFailedToCreateLogFiles
	}

	Logger.SetOutput(os.Stdout)
	if envConfig.Log.FileLogging {
		Logger.SetOutput(logFile)
	}
	return nil
}

func logWithFields(level logrus.Level, msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 && fields[0] != nil {
		Logger.WithFields(fields[0]).Log(level, msg)
	} else {
		Logger.Log(level, msg)
	}
}

func Debug(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.DebugLevel, msg, fields...)
}

func Info(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.InfoLevel, msg, fields...)
}

func Warn(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.WarnLevel, msg, fields...)
}

func Error(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.ErrorLevel, msg, fields...)
}

func Fatal(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.FatalLevel, msg, fields...)
}

func determineLogLevel(logLevel string) logrus.Level {
	switch logLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
