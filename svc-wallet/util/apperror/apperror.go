// util/apperror/apperror.go
package apperror

type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

