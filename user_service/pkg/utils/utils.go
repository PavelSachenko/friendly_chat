package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/validation"
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
	return nil, userId.(uint64)
}
