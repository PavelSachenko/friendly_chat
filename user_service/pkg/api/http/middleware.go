package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/utils"
	"github.com/pavel/user_service/pkg/validation"
	"net/http"
)

func (h *Handler) authMiddleware(ctx *gin.Context) {
	err, userId := h.auth.CheckAuthorization(utils.GetBearerToken(ctx))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, validation.IError{Field: "Bearer Token", Value: "wrong", Tag: "token"})
		return
	}

	ctx.Set("userId", userId)

	ctx.Next()
}
