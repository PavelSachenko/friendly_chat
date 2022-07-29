package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/utils"
	"github.com/pavel/user_service/pkg/validation"
	"mime/multipart"
	"net/http"
	"path/filepath"
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

type updateUserRequest struct {
	Description string `json:"description" validate:"required,max=500"`
}

func (h *Handler) updateUser(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	var updateUserRequest updateUserRequest
	requestErrors := h.validator.ValidateRequest(ctx, &updateUserRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	err, user := h.user.UpdateUser(model.UpdateUser{
		ID:          userId,
		Description: updateUserRequest.Description,
	})
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type updateAvatarRequest struct {
	Photo *multipart.FileHeader `form:"photo" binding:"required"`
}

func (h *Handler) uploadAvatarForUser(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	var updateAvatarRequest updateAvatarRequest
	err := ctx.ShouldBind(&updateAvatarRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validation.IError{Field: "photo", Tag: "required", Value: "Photo is required param"})
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("photo")
	fileExtension := filepath.Ext(fileHeader.Filename)
	if fileExtension != ".png" && fileExtension != ".jpg" && fileExtension != ".jpeg" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validation.IError{Field: "photo", Tag: "extension", Value: "File can be png,jpg,jpeg"})
		return
	}

	err = ctx.Request.ParseMultipartForm(2 << 20)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validation.IError{Field: "photo", Tag: "size", Value: "maximum photo size 2 megabytes"})
		return
	}

	err, user := h.user.UpdateAvatar(ctx, userId, file)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
