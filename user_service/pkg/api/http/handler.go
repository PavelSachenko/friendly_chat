package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/service"
	"github.com/pavel/user_service/pkg/validation"
)

type Handler struct {
	cfg       *config.Config
	log       logger.Logger
	gin       *gin.Engine
	user      service.User
	auth      service.Auth
	validator *validation.Validator
}

func InitHandler(cfg *config.Config, log logger.Logger, user service.User, auth service.Auth, gin *gin.Engine) *Handler {
	log.Printf("Init gin handler")
	return &Handler{
		cfg:       cfg,
		log:       log,
		gin:       gin,
		user:      user,
		auth:      auth,
		validator: validation.InitValidator(),
	}
}

func (h *Handler) Handle() *gin.Engine {
	h.log.Printf("Add prefix api to all handler")
	api := h.gin.Group("/api")
	//api.Use(h.cors)

	h.log.Printf("Init user handler")
	h.userHandlers(api)

	h.log.Printf("Init auth handler")
	h.authHandlers(api)

	return h.gin
}

func (h *Handler) userHandlers(api *gin.RouterGroup) {
	user := api.Group("/user")
	user.Use(h.authMiddleware)
	user.GET("", h.getUser)
	user.PUT("", h.updateUser)
	user.GET("/all", h.getFindUsers)
	user.POST("/set-avatar", h.uploadAvatarForUser)
}

func (h *Handler) authHandlers(api *gin.RouterGroup) {
	auth := api.Group("/user/auth")
	auth.POST("/sign-in", h.signIn)
	auth.POST("/sign-up", h.signUp)
	auth.POST("/logout", h.logout).Use(h.authMiddleware)
	auth.POST("/refresh", h.refreshToken).Use(h.authMiddleware)
}
