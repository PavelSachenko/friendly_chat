package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/pkg/validation"
	"net/http"
	"strings"
)

func GetBearerToken(ctx *gin.Context) string {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return token[1]
}

func GetUserIdFromContext(ctx *gin.Context) (*validation.IError, uint64) {
	userId, ok := ctx.Get("userId")
	if ok == false {
		return &validation.IError{Field: "Bearer Token", Value: "wrong", Tag: "token"}, 0
	}
	return nil, num64(userId)
}

func num64(n interface{}) uint64 {
	switch n := n.(type) {
	case int:
		return uint64(n)
	case int8:
		return uint64(n)
	case int16:
		return uint64(n)
	case int32:
		return uint64(n)
	case int64:
		return uint64(n)
	case uint:
		return uint64(n)
	case uintptr:
		return uint64(n)
	case uint8:
		return uint64(n)
	case uint16:
		return uint64(n)
	case uint32:
		return uint64(n)
	case uint64:
		return uint64(n)
	}

	return 0
}
