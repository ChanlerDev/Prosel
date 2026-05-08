package response

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Body struct {
	Data  any        `json:"data,omitempty"`
	Meta  any        `json:"meta,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, Body{Data: data})
}

func Error(c *gin.Context, status int, code string, message string, details any) {
	c.JSON(status, Body{Error: &ErrorBody{Code: code, Message: message, Details: details}})
}
