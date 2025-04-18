package models

type Response struct {
	Message string		`json:"message"`
	Status	string		`json:"status"`
	Data	interface{}	`json:"data, omitempty"`
}

func NewErrorResponse(message string, err error) Response {
	return Response {
		Message:	message,
		Status:		"error",
		Data: 		[]string{err.Error()},	
	}
}

func NewSuccessResponse(message string, data interface{}) Response {
	return Response {
		Message:	message,
		Status:		"success",
		Data:		data,
	}
}