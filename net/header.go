package net

type HeaderRow struct {
	Key, Value string
}

func NewHeaderRow(key, value string) *HeaderRow {
	return &HeaderRow{Key: key, Value: value}
}
