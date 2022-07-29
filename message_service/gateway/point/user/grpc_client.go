package user

import (
	"github.com/pavel/message_service/config"
	"github.com/pavel/message_service/gateway/point/user/pb"
	"google.golang.org/grpc"
	"log"
)

func InitUserServiceClient(cfg *config.Config) pb.UserServiceClient {
	log.Printf("Initial User grpc client")
	cc, err := grpc.Dial(cfg.Server.UserServer, grpc.WithInsecure(), grpc.WithDefaultCallOptions())
	if err != nil {
		log.Fatalf("Can not connect to User service: %v", err)
	}

	return pb.NewUserServiceClient(cc)
}
