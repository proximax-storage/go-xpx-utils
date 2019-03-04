package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
)

type Logger struct {
	*zap.Logger
	operationKey  string
	encoderConfig zapcore.EncoderConfig
	cores         []zapcore.Core
}

func NewLoggerFromZapConfig(config zap.Config) (*Logger, error) {
	zapLogger, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse logger config")
	}

	logger := &Logger{
		Logger:        zapLogger,
		encoderConfig: config.EncoderConfig,
	}

	return logger, nil
}

func (ref *Logger) SetLogsTransfer(identifierField, identifierValue string, transfer *Transfer) error {
	if transfer == nil {
		return ErrNilTransfer
	}

	var fieldInterceptor fieldsInterceptorFunc

	if len(identifierField) != 0 && len(identifierValue) != 0 {
		fieldInterceptor = fieldsInterceptorFunc(func(fields []zapcore.Field) []zapcore.Field {
			return append(fields, zap.String(identifierField, identifierValue))
		})
	}

	ref.cores = append(ref.cores, newInterceptedCore(
		ref.Core(),
		zapcore.NewJSONEncoder(ref.encoderConfig),
		fieldInterceptor,
		logStringInterceptorFunc(func(log string) string {
			transfer.Pool() <- log

			return log
		}),
	))

	return nil
}

func (ref *Logger) SetOperationKey(key string) error {
	if len(key) == 0 {
		return ErrBlankOperationKey
	}

	ref.operationKey = key

	return nil
}

func (ref *Logger) addContextField(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx != nil && ctx.Value(ref.operationKey) != nil {
		fields = append(
			fields,
			zap.Any(ref.operationKey, ctx.Value(ref.operationKey)),
		)
	}

	return fields
}

func (ref *Logger) log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if ce := ref.Check(lvl, msg); ce != nil {
		if len(ref.cores) != 0 {
			for _, core := range ref.cores {
				ce.AddCore(ce.Entry, core)
			}
		}

		ce.Write(fields...)
	}
}

func (ref *Logger) Debug(msg string, fields ...zap.Field) {
	ref.log(zapcore.DebugLevel, msg, fields...)
}

func (ref *Logger) Info(msg string, fields ...zap.Field) {
	ref.log(zapcore.InfoLevel, msg, fields...)
}

func (ref *Logger) Error(msg string, fields ...zap.Field) {
	ref.log(zapcore.ErrorLevel, msg, fields...)
}

func (ref *Logger) Warn(msg string, fields ...zap.Field) {
	ref.log(zapcore.WarnLevel, msg, fields...)
}

func (ref *Logger) DPanic(msg string, fields ...zap.Field) {
	ref.log(zapcore.DPanicLevel, msg, fields...)
}

func (ref *Logger) Panic(msg string, fields ...zap.Field) {
	ref.log(zapcore.PanicLevel, msg, fields...)
}

func (ref *Logger) Fatal(msg string, fields ...zap.Field) {
	ref.log(zapcore.FatalLevel, msg, fields...)
}

func (ref *Logger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsDebugEnabled() {
		ref.Debug(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsInfoEnabled() {
		ref.Info(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsWarnEnabled() {
		ref.Warn(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsErrorEnabled() {
		ref.Error(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) DPanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsDPanicEnabled() {
		ref.DPanic(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsPanicEnabled() {
		ref.Panic(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	if ref.IsFatalEnabled() {
		ref.Fatal(msg, ref.addContextField(ctx, fields)...)
	}
}

func (ref *Logger) IsDebugEnabled() bool {
	return ref.Core().Enabled(zap.DebugLevel)
}

func (ref *Logger) IsInfoEnabled() bool {
	return ref.Core().Enabled(zap.InfoLevel)
}

func (ref *Logger) IsWarnEnabled() bool {
	return ref.Core().Enabled(zap.WarnLevel)
}

func (ref *Logger) IsErrorEnabled() bool {
	return ref.Core().Enabled(zap.ErrorLevel)
}

func (ref *Logger) IsDPanicEnabled() bool {
	return ref.Core().Enabled(zap.DPanicLevel)
}

func (ref *Logger) IsPanicEnabled() bool {
	return ref.Core().Enabled(zap.PanicLevel)
}

func (ref *Logger) IsFatalEnabled() bool {
	return ref.Core().Enabled(zap.FatalLevel)
}
