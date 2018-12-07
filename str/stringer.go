package str

import (
	"fmt"
	"strings"
)

type FieldPattern string

const (
	StringPattern  FieldPattern = "s"
	IntPattern     FieldPattern = "d"
	BooleanPattern FieldPattern = "t"
	FloatPattern   FieldPattern = "g"
	ValuePattern   FieldPattern = "v"
)

type field struct {
	fieldName string
	pattern   FieldPattern
	value     interface{}
}

// NewField returns field which is used for StructToString() function
func NewField(fieldName string, pattern FieldPattern, val interface{}) *field {
	return &field{
		fieldName: fieldName,
		pattern:   pattern,
		value:     val,
	}
}

func (f *field) String() string {
	return fmt.Sprintf(f.fieldName+"=%"+string(f.pattern), f.value)
}

// StructToString building string by fields from object you provided
// Use it within implementation fmt.Stringer interface
func StructToString(structName string, fields ...*field) string {
	if len(fields) == 0 {
		return ""
	}

	values := make([]string, len(fields))

	for idx, field := range fields {
		values[idx] = field.String()
	}

	return fmt.Sprintf("%s [%s]", structName, strings.Join(values, ", "))
}
