package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/service"
)

type Handler struct {
	log  logger.Logger
	gin  *gin.Engine
	user service.User
	auth service.Auth
}

func InitHandler(log logger.Logger, user service.User, auth service.Auth) *Handler {
	log.Printf("Init gin handler")
	return &Handler{
		log:  log,
		gin:  gin.New(),
		user: user,
		auth: auth,
	}
}

func (h *Handler) Handle() *gin.Engine {
	h.log.Printf("Add prefix api to all handlers")
	api := h.gin.Group("/api/user")

	h.log.Printf("Init user handlers")
	h.userHandlers(api)

	h.log.Printf("Init auth handlers")
	h.authHandlers(api)

	return h.gin
}

func (h *Handler) userHandlers(api *gin.RouterGroup) {
	api.Use(h.authMiddleware).GET("/", h.getUser)
}

func (h *Handler) authHandlers(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	auth.POST("/sig-in", h.sighIn)
	auth.POST("/sig-up", h.sighUp).Use(h.authMiddleware)
	auth.POST("/logout", h.logout).Use(h.authMiddleware)
	auth.POST("/refresh", h.refreshToken).Use(h.authMiddleware)
}
