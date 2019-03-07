package logger

import (
	"go.uber.org/zap"
)

type ZapFieldCreator func() zap.Field
type ZapFieldsCreator func() []zap.Field

func ZapField(creator ZapFieldCreator) zap.Field {
	return creator()
}

func ZapFields(creator ZapFieldsCreator) []zap.Field {
	return creator()
}
