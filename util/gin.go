package util

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const (
	StatusOK          = 200
	StatusParamsErr   = 1001
	StatusInternalErr = 1002
)

type Response struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"traceId"`
}

var StatusMessage = map[int32]string{
	200:  "success",
	1001: "参数错误",
	1002: "系统出错",
}

func GinErrResponse(ctx *gin.Context, err error) {
	traceID := GetGinCtxTraceID(ctx)
	var validationError validator.ValidationErrors
	if errors.As(err, &validationError) {
		ctx.JSON(http.StatusBadRequest, &Response{
			Code:    StatusParamsErr,
			Message: StatusMessage[StatusParamsErr],
			Data:    err,
			TraceID: traceID,
		})
		return
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		ctx.JSON(http.StatusBadRequest, &Response{
			Code:    StatusParamsErr,
			Message: StatusMessage[StatusParamsErr],
			Data:    err,
			TraceID: traceID,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, &Response{
		Code:    StatusInternalErr,
		Message: StatusMessage[StatusInternalErr],
		Data:    err,
		TraceID: traceID,
	})
	return
}

func GinSuccessResponse(ctx *gin.Context, res interface{}) {
	traceID := GetGinCtxTraceID(ctx)
	ctx.JSON(http.StatusOK, &Response{
		Code:    StatusOK,
		Message: StatusMessage[StatusOK],
		Data:    res,
		TraceID: traceID,
	})
}

func GetGinCtxTraceID(ctx *gin.Context) string {
	var traceID string
	if val, ok := ctx.Get(CtxTraceID); ok {
		traceID = val.(string)
	}
	return traceID
}
