package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pavel/user_service/config"
	api "github.com/pavel/user_service/pkg/api/http"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/pb"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/pavel/user_service/pkg/service"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	log := logger.GetLogger()
	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Fatalf("Error init config: ERROR: %s", err.Error())
	}
	err, postgres := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("Error init postgres db: CONFIG: %v ERROR: %s", cfg.DB, err.Error())
	}
	err, redis := db.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Error init postgres db: CONFIG: %v ERROR: %s", cfg.DB, err.Error())
	}
	userRepo := repository.InitUserPostgres(postgres, log)
	authRepo := repository.InitAuthRedis(redis, postgres, log)

	userService := service.InitUserService(userRepo, cfg, log)
	authService := service.InitAuthService(authRepo, cfg)
	handlerOrigin := gin.New()
	handlerOrigin.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",
			"http://localhost:1000",
			"http://localhost:10000",
			"http://localhost:10001",
			"http://localhost:10002",
			"http://localhost:3000",
			"http://localhost:3001",
		},
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders: []string{
			"Authorization",
			"Access-Control-Allow-Headers",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Content-Type",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Sec-WebSocket-Protocol",
			"*"},
		AllowCredentials: true,
		AllowWildcard:    false,
		AllowWebSockets:  true,
		MaxAge:           84000,
	}))

	handlerOrigin.Use(func(ctx *gin.Context) {
		ctx.Set("aws_session", connectToAws(cfg))
		ctx.Next()
	})
	handlerOrigin.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.Next()
	})
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, pb.InitGRPCServer(log, userService, authService))
	lis, err := net.Listen("tcp", cfg.Server.GRPCAddress)
	if err != nil {
		log.Fatalf("Error init grpc server: %s", err.Error())
	}
	wg.Add(1)
	go serveGRPC(lis, grpcServer)
	api.InitHandler(cfg, *log, userService, authService, handlerOrigin).Handle()
	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, handlerOrigin)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
	wg.Wait()
}

func serveGRPC(lis net.Listener, grpcServer *grpc.Server) {
	defer wg.Done()
	grpcServer.Serve(lis)
}

func connectToAws(cfg *config.Config) *session.Session {
	accessKeyID := cfg.Aws.AccessKey
	secretAccessKey := cfg.Aws.SecretKey
	myRegion := cfg.Aws.Region
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(myRegion),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}
