package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/message_service/gateway/point/user/pb"
	"github.com/pavel/message_service/pkg/utils"
	"github.com/pavel/message_service/pkg/validation"
	"golang.org/x/net/context"
	"net/http"
)

type AuthMiddleware struct {
	userSvc pb.UserServiceClient
}

func InitAuthMiddleware(userSvc pb.UserServiceClient) AuthMiddleware {
	return AuthMiddleware{
		userSvc: userSvc,
	}
}

func (c *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	res, err := c.userSvc.CheckToken(context.Background(), &pb.CheckTokenRequest{
		Token: utils.GetBearerToken(ctx),
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, validation.IError{Field: "Bearer Token", Value: "wrong", Tag: "token"})
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}
