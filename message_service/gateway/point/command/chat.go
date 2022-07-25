package command

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/gateway/point/user/pb"
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/utils"
	"net/http"
)

type createPrivateChatRequest struct {
	UserId   int    `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required,min=1"`
}

func (h Handler) createPrivateChat(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	var createPrivateChatRequest createPrivateChatRequest
	requestErrors := h.validator.ValidateRequest(ctx, &createPrivateChatRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	res, err := h.userSvc.GetUser(ctx, &pb.GetUserRequest{
		UserId: int32(userId),
	})

	if err != nil || res.Status >= 400 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, res.Err)
	}

	chatId, err := h.chat.CreateBetweenUsers([]model.UserTitleChats{
		{UserId: int(userId), Title: createPrivateChatRequest.Username},
		{UserId: createPrivateChatRequest.UserId, Title: res.User.Username},
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
	}

	ctx.JSON(200, chatId)
}

type getAllChatRequest struct {
	Limit  int `form:"limit,default=20" json:"limit"`
	Offset int `form:"offset,default=0" json:"offset"`
}

func (h Handler) getAllChats(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}

	var getAllChatRequest getAllChatRequest
	requestErrors := h.validator.ValidateRequest(ctx, &getAllChatRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	chats, err := h.chat.GetAllForUser(model.FilterChat{
		UserId: int(userId),
		Limit:  getAllChatRequest.Limit,
		Offset: getAllChatRequest.Offset,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, chats)
}
