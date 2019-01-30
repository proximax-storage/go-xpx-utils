package net

type IdentifiableError struct {
	ErrorId string        `json:"error_id"`
	Message string        `json:"message"`
	Args    []interface{} `json:"args,omitempty"`
}

func (ref *IdentifiableError) Error() string {
	return ref.Message
}
