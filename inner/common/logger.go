package common

import (
	"context"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ключ для получения requestId из контекста
var ridKey = requestid.ConfigDefault.ContextKey.(string)

// Logger структура логгера
type Logger struct {
	*zap.Logger
}

// NewLogger функция-конструктор логгера
func NewLogger(cfg Config) *Logger {
	zapEncoderCfg := zapcore.EncoderConfig{
		TimeKey:          "timestamp",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000"),
		EncodeDuration:   zapcore.MillisDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "  ",
	}

	zapCfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(parseLogLevel(cfg.LogLevel)),
		Development: cfg.LogDevelopMode,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		// пишем записи в формате JSON
		Encoding:      "json",
		EncoderConfig: zapEncoderCfg,
		// логируем сообщения и ошибки в консоль
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	logger := zap.Must(zapCfg.Build())
	logger.Info("logger construction succeeded")
	created := &Logger{logger}
	created.setNewFiberZapLogger()
	return created
}

// setNewFiberZapLogger устанавливает логгер для fiber
func (l *Logger) setNewFiberZapLogger() {
	fiberZapLogger := fiberzap.NewLogger(fiberzap.LoggerConfig{
		SetLogger: l.Logger,
	})
	log.SetLogger(fiberZapLogger)
}

// parseLogLevel парсит уровень логирования из строки в zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO":
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "panic", "PANIC":
		return zapcore.PanicLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// функция логирования с добавлением requestId
func (l *Logger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	// получаем requestId из контекста
	var rid string
	if v := ctx.Value(ridKey); v != nil {
		rid = v.(string)
	}
	// логируем
	l.Debug(msg, append(fields, zap.String("requestid", rid))...)
	// добавляем для логирования поле с requestId
	fields = append(fields, zap.String(ridKey, rid))
	// вызываем метод логгера
	l.Debug(msg, fields...)
}

// функция логирования с добавлением requestId
func (l *Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	// получаем requestId из контекста
	var rid string
	if v := ctx.Value(ridKey); v != nil {
		rid = v.(string)
	}
	// логируем
	l.Error(msg, append(fields, zap.String("requestid", rid))...)
	// добавляем для логирования поле с requestId
	fields = append(fields, zap.String(ridKey, rid))
	// вызываем метод логгера
	l.Error(msg, fields...)
}
