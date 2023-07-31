package logger

import (
	"os"

	"github.com/assbomber/myzone/configs"
	"github.com/assbomber/myzone/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func (l *Logger) Info(message string) {
	l.log.Info(message)
}

func (l *Logger) Debug(message string) {
	l.log.Debug(message)
}

func (l *Logger) Error(message string, err error) {
	l.log.Error(message, zap.Error(err))
}

func (l *Logger) Fatal(message string) {
	l.log.Fatal(message)
}

// Initializes the logger
func InitLogger() *Logger {
	logger := Logger{}
	if configs.GetString("RUNTIME_ENV") == constants.Environments.PRODUCTION {
		logger.log = initializeFileLogger()
	} else {
		logger.log = initializeConsoleLogger()
	}
	return &logger
}

// Initializes the logger for console logging. Ideal all environments other than production
func initializeConsoleLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err.Error())
	}
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

// Initializes the logger for file logging. Ideal for all environments other than development
func initializeFileLogger() *zap.Logger {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(conf)

	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err.Error())
	}

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}
