package logger

import (
	"context"

	"go.uber.org/zap"
)

var staticLogger = &Logger{
	Logger: zap.NewNop(),
}

func SetStaticLogger(logger *Logger) error {
	if logger == nil {
		return ErrNilLogger
	}

	staticLogger = logger

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	staticLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	staticLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	staticLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	staticLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	staticLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	staticLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	staticLogger.Fatal(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.DebugCtx(ctx, msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.InfoCtx(ctx, msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.WarnCtx(ctx, msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.ErrorCtx(ctx, msg, fields...)
}

func DPanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.DPanicCtx(ctx, msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.PanicCtx(ctx, msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	staticLogger.FatalCtx(ctx, msg, fields...)
}

func IsDebugEnabled() bool {
	return staticLogger.IsDebugEnabled()
}

func IsInfoEnabled() bool {
	return staticLogger.IsInfoEnabled()
}

func IsWarnEnabled() bool {
	return staticLogger.IsWarnEnabled()
}

func IsErrorEnabled() bool {
	return staticLogger.IsErrorEnabled()
}

func IsDPanicEnabled() bool {
	return staticLogger.IsDPanicEnabled()
}

func IsPanicEnabled() bool {
	return staticLogger.IsPanicEnabled()
}

func IsFatalEnabled() bool {
	return staticLogger.IsFatalEnabled()
}

func Sync() error {
	return staticLogger.Sync()
}
