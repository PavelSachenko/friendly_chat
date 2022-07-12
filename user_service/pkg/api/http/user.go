package http

import "github.com/gin-gonic/gin"

func (h *Handler) getUser(ctx *gin.Context) {
	ctx.JSON(200, "This my user")
}
