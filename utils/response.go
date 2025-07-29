package utils

type ErrorResponse struct {
	Errors []map[string]string `json:"errors"`
}

type ValidationErrorResponse ErrorResponse

func NewValidationErrorResponse(errors []map[string]string) ValidationErrorResponse {
	return ValidationErrorResponse{Errors: errors}
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Errors: []map[string]string{
			{"field": "internal", "message": message},
		},
	}
}
