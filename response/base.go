package response

import (
	"github.com/tuyenlqvnp/sign-service-api/bean"
)

type Base struct {
	Status       int         `json:"status"`
	Message      string      `json:"message"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}

func (e *Base) SetStatus(key string, err error) {
	e.Status = bean.CodeMessage[key].Code
	e.Message = bean.CodeMessage[key].Message
	if err != nil {
		e.ErrorMessage = err.Error()
	}
}
