package net

type IdentifiableError struct {
	ErrorId       uint16 `json:"error_id"`
	ErrorStringId string `json:"error_string_id"`
	Message       string `json:"message"`
}

func (ref *IdentifiableError) Error() string {
	return ref.Message
}
