package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/model"
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
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type findUserRequest struct {
	Limit    int    `form:"limit"  validate:"numeric"`
	Offset   int    `form:"offset" validate:"numeric"`
	Username string `form:"username"`
}

func (h *Handler) getFindUsers(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	var findUserRequest findUserRequest
	requestErrors := h.validator.ValidateRequest(ctx, &findUserRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	err, users := h.user.FindUser(model.UserFilter{
		OwnerUserId: userId,
		Limit:       findUserRequest.Limit,
		Offset:      findUserRequest.Offset,
		Username:    findUserRequest.Username,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)

}
