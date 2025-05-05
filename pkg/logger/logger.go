package logger

import (
	"context"
	"log"
	"sync"

	"app/internal/config"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once     sync.Once
	instance *Logger
)

type Logger struct {
	core *zap.SugaredLogger
	ctx  context.Context
}

func (l *Logger) With(key string, value interface{}) *Logger {
	return &Logger{
		core: l.core.With(key, value),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		core: l.core.With("error", err.Error()),
	}
}

func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		core: l.core.With("service", service),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		ctx = context.WithValue(context.Background(), "requestId", uuid.NewString())
	} else if requestId := ctx.Value("requestId"); requestId == nil || requestId == "" {
		requestId = uuid.NewString()
		ctx = context.WithValue(ctx, "requestId", requestId)
	}

	return &Logger{
		core: l.core.With("requestId", ctx.Value("requestId")),
		ctx:  ctx,
	}
}

func (l *Logger) Debug(msg string) {
	l.core.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.core.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.core.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.core.Error(msg)
}

func (l *Logger) Panic(msg string) {
	l.core.Panic(msg)
}

func New() *Logger {
	once.Do(func() {
		cfg := config.New().Logger
		zCfg := zap.NewProductionConfig()
		zCfg.EncoderConfig.LevelKey = "logLevel"
		zCfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

		lvlM := map[string]zapcore.Level{
			"debug": zap.DebugLevel,
			"info":  zap.InfoLevel,
			"warn":  zap.WarnLevel,
			"error": zap.ErrorLevel,
		}

		if lvl, ok := lvlM[cfg.Level]; ok {
			zCfg.Level = zap.NewAtomicLevelAt(lvl)
		} else {
			zCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			log.Printf("invalid log level: %s use debug level as default\n", cfg.Level)
		}

		logger, _ := zCfg.Build(zap.AddCallerSkip(1))
		defer logger.Sync()

		instance = &Logger{
			core: logger.Sugar(),
		}

		instance.Info("logger initialized")
	})

	return instance
}
