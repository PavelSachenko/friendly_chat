package pb

import (
	"context"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type grpcUserServer struct {
	log  *logger.Logger
	user service.User
	auth service.Auth
	UserServiceServer
}

func InitGRPCServer(log *logger.Logger, user service.User, auth service.Auth) *grpcUserServer {
	return &grpcUserServer{
		log:  log,
		user: user,
		auth: auth,
	}
}

func (g grpcUserServer) CheckToken(ctx context.Context, request *CheckTokenRequest) (*CheckTokenResponse, error) {
	err, userId := g.auth.CheckAuthorization(request.Token)
	if err != nil {
		return &CheckTokenResponse{
			Status: http.StatusUnauthorized,
			Err:    err.Error(),
		}, nil
	}
	return &CheckTokenResponse{
		Status: http.StatusOK,
		Err:    "",
		UserId: int32(userId),
	}, nil
}

func (g grpcUserServer) GetUser(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	err, user := g.user.GetUser(uint64(request.UserId))
	if err != nil {
		return &GetUserResponse{
			Status: http.StatusUnprocessableEntity,
			Err:    err.Error(),
		}, nil
	}

	return &GetUserResponse{
		Status: http.StatusOK,
		User: &User{
			Id:          int32(user.ID),
			Username:    user.Username,
			Description: user.Description,
			Avatar:      user.Avatar,
			CreatedAt:   user.CreatedAt.Unix(),
		},
	}, nil
}
