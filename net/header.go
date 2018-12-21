package net

type HeaderRaw struct {
	key, value string
}

func NewHeaderRow(key, value string) *HeaderRaw {
	return &HeaderRaw{key: key, value: value}
}
