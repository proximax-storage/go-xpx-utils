package net

type IdentifiableError struct {
	ErrorId string
	Message string
}

func (ref *IdentifiableError) Error() string {
	return ref.Message
}
