package command

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/config"
	"github.com/pavel/message_service/gateway/point/user"
	"github.com/pavel/message_service/gateway/point/user/pb"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/service"
	"github.com/pavel/message_service/pkg/validation"
)

type Handler struct {
	log       *logger.Logger
	cfg       *config.Config
	userSvc   pb.UserServiceClient
	chat      service.Chat
	message   service.Message
	validator *validation.Validator
}

func InitHandler(
	r *gin.RouterGroup,
	log *logger.Logger,
	cfg *config.Config,
	userSvc pb.UserServiceClient,
	chat service.Chat,
	message service.Message,
) {
	h := Handler{
		log:       log,
		cfg:       cfg,
		userSvc:   userSvc,
		chat:      chat,
		message:   message,
		validator: validation.InitValidator(),
	}

	h.registerHandlers(r)
}

func (h Handler) registerHandlers(r *gin.RouterGroup) {
	auth := user.InitAuthMiddleware(h.userSvc)
	chat := r.Group("/chat")
	chat.Use(auth.AuthRequired)
	message := chat.Group("/")

	chat.POST("/private", h.createPrivateChat)
	chat.POST("/group", h.createPrivateChat)
	chat.GET("/all", h.getAllChats)

	message.POST("/message", h.sendMessage)
	message.GET("/messages", h.getAll)
}
