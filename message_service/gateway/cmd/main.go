package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pavel/message_service/config"
	"github.com/pavel/message_service/gateway/point/command"
	"github.com/pavel/message_service/gateway/point/user"
	"github.com/pavel/message_service/pkg/db"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/repository"
	"github.com/pavel/message_service/pkg/service"
)

func main() {
	log := logger.GetLogger()
	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Warnf("Error init config: ERROR: %s", err.Error())
	}
	err, postgre := db.InitPostgres(cfg, db.InitPostgresQueryBuilder())
	if err != nil {
		log.Warnf("Error init postgre: ERROR: %s", err.Error())
	}

	chatRepo := repository.InitChatPostgreSQL(*postgre, *log)
	chatService := service.InitChatService(chatRepo)

	messageRepo := repository.InitMessagePostgreSQL(postgre, log)
	messageService := service.InitMessageService(messageRepo)

	userSvc := user.InitUserServiceClient(cfg)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:1000", "http://localhost:10000", "http://localhost:10001", "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Sec-WebSocket-Protocol", "*"},
		AllowCredentials: true,
		AllowWildcard:    false,
		AllowWebSockets:  true,

		MaxAge: 84000,
	}))
	r.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Content-Type", "application/json")
	})
	api := r.Group("/api/message_service")
	command.InitHandler(api, log, cfg, userSvc, chatService, messageService)

	log.Fatal(r.Run(cfg.Server.GatewayAddress))
}
