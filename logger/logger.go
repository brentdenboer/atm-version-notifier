package logger

import (
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents the logging level
type Level string

const (
	DEBUG   Level = "debug"
	INFO    Level = "info"
	WARNING Level = "warn"
	ERROR   Level = "error"
	FATAL   Level = "fatal"
)

var (
	defaultLogger *zap.SugaredLogger
	levelAtom     zap.AtomicLevel // Atomic level for dynamic level changes
)

func init() {
	levelAtom = zap.NewAtomicLevel()
	levelAtom.SetLevel(zapcore.InfoLevel)

	// Start with a default stdout logger using atomic level
	SetOutput(os.Stdout)
}

// New creates a new logger instance with the given writers
func New(writers ...io.Writer) *zap.SugaredLogger {
	var cores []zapcore.Core

	// Default encoder config for human-readable output
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create a core for each writer
	for _, w := range writers {
		var encoder zapcore.Encoder
		if _, isFile := w.(*os.File); isFile {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(w),
			levelAtom,
		)
		cores = append(cores, core)
	}

	// Create a logger with all cores
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger.Sugar()
}

// SetOutput sets the output for the default logger
func SetOutput(writers ...io.Writer) {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	defaultLogger = New(writers...)
}

// SetLevel sets the minimum logging level
func SetLevel(level Level) {
	var zapLevel zapcore.Level
	switch strings.ToLower(string(level)) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	levelAtom.SetLevel(zapLevel)
}

// Package-level logging functions that use the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Sync ensures all buffered logs are written
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}
