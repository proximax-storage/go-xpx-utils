package net

type IdentifiableError struct {
	ErrorId string        `json:"error_id,omitempty"`
	Message string        `json:"message,omitempty"`
	Args    []interface{} `json:"args,omitempty"`
}

func (ref *IdentifiableError) Error() string {
	return ref.Message
}
