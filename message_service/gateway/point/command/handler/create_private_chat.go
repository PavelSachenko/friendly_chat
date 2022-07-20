package handler

import "github.com/gin-gonic/gin"

func CreatePrivateChat(ctx *gin.Context) {
	ctx.JSON(200, "Hello world")
}
