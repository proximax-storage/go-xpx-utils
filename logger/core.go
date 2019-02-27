package logger

import (
	"go.uber.org/zap/zapcore"
)

type fieldsInterceptor interface {
	intercept(fields []zapcore.Field) []zapcore.Field
}

type fieldsInterceptorFunc func(fields []zapcore.Field) []zapcore.Field

func (ref fieldsInterceptorFunc) intercept(fields []zapcore.Field) []zapcore.Field {
	return ref(fields)
}

type logStringInterceptor interface {
	intercept(log string) string
}

type logStringInterceptorFunc func(log string) string

func (ref logStringInterceptorFunc) intercept(log string) string {
	return ref(log)
}

type interceptedCore struct {
	zapcore.Core
	encoder              zapcore.Encoder
	fieldsInterceptor    fieldsInterceptor
	logStringInterceptor logStringInterceptor
}

func newInterceptedCore(core zapcore.Core, enc zapcore.Encoder, fieldsInterceptor fieldsInterceptor, logStringInterceptor logStringInterceptor) zapcore.Core {
	return &interceptedCore{
		Core:                 core,
		encoder:              enc,
		fieldsInterceptor:    fieldsInterceptor,
		logStringInterceptor: logStringInterceptor,
	}
}

func (ref *interceptedCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if ref.fieldsInterceptor != nil {
		fields = ref.fieldsInterceptor.intercept(fields)
	}

	if ref.logStringInterceptor != nil {
		buf, err := ref.encoder.EncodeEntry(entry, fields)
		if err != nil {
			return err
		}

		ref.logStringInterceptor.intercept(buf.String())

		buf.Free()
	}

	return nil
}
