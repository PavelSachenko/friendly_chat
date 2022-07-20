package command

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/config"
	"github.com/pavel/message_service/gateway/point/command/handler"
)

func RegisterHandlers(r *gin.RouterGroup, cfg *config.Config) {
	chat := r.Group("/chat")
	chat.POST("/create-private", createPrivateChat)
	chat.POST("/create-group", createGroupChat)
}

func createPrivateChat(ctx *gin.Context) {
	handler.CreatePrivateChat(ctx)
}

func createGroupChat(ctx *gin.Context) {
	handler.CreatePrivateChat(ctx)
}
