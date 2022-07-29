package command

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/utils"
	"net/http"
)

type sendMessageRequest struct {
	ChatID int    `json:"chat_id" validate:"required"`
	Body   string `json:"body" validate:"required,min=1,max=5000"`
}

func (h Handler) sendMessage(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	var sendMessageRequest sendMessageRequest
	requestErrors := h.validator.ValidateRequest(ctx, &sendMessageRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	messageId, err := h.message.Send(int(userId), sendMessageRequest.ChatID, sendMessageRequest.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, messageId)
}

type getAllMessageRequest struct {
	ChatID int `form:"chat_id" json:"chat_id" validate:"required"`
	Limit  int `form:"limit,default=20" json:"limit"`
	Offset int `form:"offset,default=0" json:"offset"`
}

func (h Handler) getAll(ctx *gin.Context) {
	authError, userId := utils.GetUserIdFromContext(ctx)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, authError)
		return
	}
	var getAllMessageRequest getAllMessageRequest
	requestErrors := h.validator.ValidateRequest(ctx, &getAllMessageRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	messages, err := h.message.GetAllForChat(model.FilterMessage{
		ChatId: getAllMessageRequest.ChatID,
		UserId: int(userId),
		Offset: getAllMessageRequest.Offset,
		Limit:  getAllMessageRequest.Limit,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, messages)
}
