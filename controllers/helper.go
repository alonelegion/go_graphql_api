package controllers

import "github.com/gin-gonic/gin"

// HTTP Response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HTTPResponse(context *gin.Context, httpCode int, message string, data interface{}) {
	context.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    data,
	})
	return
}
