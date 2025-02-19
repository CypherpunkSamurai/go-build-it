// logger.go (normal version written by me)
package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger is the logger that is used
var zapLogger *zap.Logger

// InitZapLogger creates a new zap logger
func InitZapLogger() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// copied configs from:
	// - https://betterstack.com/community/guides/logging/go/zap/
	// - https://last9.io/blog/zap-logger/
	var config zap.Config
	if IsProd() {
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
			InitialFields: map[string]interface{}{
				"pid": os.Getpid(),
			},
		}
	} else {
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	// Must be called before using the logger
	zapLogger = zap.Must(config.Build())
}

// Logger returns a zap logger that prints to stdout
// Example:
// ```go
// utils.Logger().Println("Hello World")
// ```
func Logger() *log.Logger {
	if zapLogger == nil {
		InitZapLogger()
	}
	return zap.NewStdLog(zapLogger)
}
