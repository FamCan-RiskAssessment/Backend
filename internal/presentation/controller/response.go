package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SingleMessageResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MultipleMessageResponse struct {
	Status   int                          `json:"status"`
	Messages map[string]map[string]string `json:"messages"`
	Data     interface{}                  `json:"data"`
}

func Response[T string | map[string]map[string]string](ctx *gin.Context, statusCode int, message T, data interface{}) {
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		panic(fmt.Errorf("invalid status code: %d", statusCode))
	}

	switch msg := any(message).(type) {
	case string:
		if msg == "" {
			msg = statusText
		}
		ctx.JSON(statusCode, SingleMessageResponse{
			Status:  statusCode,
			Message: msg,
			Data:    data,
		})
	case map[string]map[string]string:
		ctx.JSON(statusCode, MultipleMessageResponse{
			Status:   statusCode,
			Messages: msg,
			Data:     data,
		})
	}
}
