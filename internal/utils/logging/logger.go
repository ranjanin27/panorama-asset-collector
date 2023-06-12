// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package logging

import (
	"context"
	"fmt"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/contextutilities"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Print(args ...interface{})
	Printf(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})

	Debug(args ...interface{})
	Debugf(template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})

	WithFields(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger

	Close()
}

// PFHLogger is a wrapper over a logging library
type PFHLogger struct {
	zap *zap.SugaredLogger
}

var (
	logger Logger
)

//nolint:gochecknoinits // This can be ignored
func init() {
	// Create zap config
	logConfig := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			FunctionKey: "function",
		},
	}

	// set log level in config if provided
	if len(configs.GetLogLevel()) > 0 {
		err := logConfig.Level.UnmarshalText([]byte(configs.GetLogLevel()))
		if err != nil {
			panic(err)
		}
	}

	// Build the logger using the config and add the caller skip option
	// Caller skip will help skip 1 level in the call stack to print the caller instead of the logger wrapper
	zaplogger, err := logConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger = &PFHLogger{zaplogger.Sugar()}
}

func GetLogger() Logger {
	return logger
}

func (h *PFHLogger) Info(args ...interface{}) {
	h.zap.Info(args...)
}

func (h *PFHLogger) Infof(template string, args ...interface{}) {
	h.zap.Infof(template, args...)
}

// Using Info for Printf. Printf is required for Kafka logging
func (h *PFHLogger) Print(args ...interface{}) {
	h.zap.Debug(args...)
}

func (h *PFHLogger) Printf(template string, args ...interface{}) {
	h.zap.Debugf(template, args...)
}

func (h *PFHLogger) Error(args ...interface{}) {
	h.zap.Error(args...)
}

func (h *PFHLogger) Errorf(template string, args ...interface{}) {
	h.zap.Errorf(template, args...)
}

func (h *PFHLogger) Warn(args ...interface{}) {
	h.zap.Warn(args...)
}

func (h *PFHLogger) Warnf(template string, args ...interface{}) {
	h.zap.Warnf(template, args...)
}

func (h *PFHLogger) Debug(args ...interface{}) {
	h.zap.Debug(args...)
}

func (h *PFHLogger) Debugf(template string, args ...interface{}) {
	h.zap.Debugf(template, args...)
}

func (h *PFHLogger) Fatal(args ...interface{}) {
	h.zap.Fatal(args...)
}

func (h *PFHLogger) Fatalf(template string, args ...interface{}) {
	h.zap.Fatalf(template, args...)
}

func (h *PFHLogger) WithFields(args ...interface{}) Logger {
	return &PFHLogger{zap: h.zap.With(args...)}
}

// WithContext extracts the corresponding ID information and appends it to the logger.
func (h *PFHLogger) WithContext(ctx context.Context) Logger {
	trace := contextutilities.ExtractTraceDetails(ctx)
	if len(trace.TraceID) > 0 && len(trace.SpanID) > 0 {
		return h.WithFields(contextutilities.LoggerTraceID, trace.TraceID,
			contextutilities.LoggerSpanID, trace.SpanID)
	}
	return h
}

// Close implemented here as Sync returns an error, so can't be deferred
func (h *PFHLogger) Close() {
	if err := h.zap.Sync(); err != nil {
		fmt.Println("Sync on zap logger failed.", err)
	}
}
