package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/utils"
	"net/http"
)

func (h *Handler) getUser(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	err, user := h.user.GetUser(userId)

	if err != nil {
		ctx.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}
