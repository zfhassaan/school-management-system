package dto

// Envelope is the uniform API response envelope.
type Response[T any] struct {
	Success bool       `json:"success"`
	Data    T          `json:"data,omitempty"`
	Error   *ErrorDTO  `json:"error,omitempty"`
}

type ErrorDTO struct {
	Code          string            `json:"code"`
	Message       string            `json:"message"`
	Fields        map[string]string `json:"fields,omitempty"`
	CorrelationID string            `json:"correlationId,omitempty"`
}

func Ok[T any](data T) Response[T] {
	return Response[T]{Success: true, Data: data}
}

func Fail[T any](code, message string, fields map[string]string, correlationID string) Response[T] {
	return Response[T]{
		Success: false,
		Error: &ErrorDTO{
			Code:          code,
			Message:       message,
			Fields:        fields,
			CorrelationID: correlationID,
		},
	}
}

