package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pavel/user_service/config"
	api "github.com/pavel/user_service/pkg/api/http"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

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

	userService := service.InitUserService(userRepo)
	authService := service.InitAuthService(authRepo, cfg)
	test := gin.New()
	test.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1000", "http://localhost:10000", "http://localhost:10001"},
		AllowMethods: []string{"*"},
		//AllowHeaders:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Sec-WebSocket-Protocol", "*"},
		AllowCredentials: true,
		AllowWildcard:    false,
		AllowWebSockets:  true,

		MaxAge: 84000,
		//AllowAllOrigins: true,
		//AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return true
		//},
	}))
	api.InitHandler(*log, userService, authService, test).Handle()
	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, test)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}

//func cors(ctx *gin.Context) {
//	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//	ctx.Writer.Header().Set("Content-Type", "application/json")
//	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
//	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
//	ctx.Writer.Header().Set("Access-Control-Request-Method", "*")
//	ctx.Writer.Header().Set("Access-Control-Request-Headers", "*")
//	ctx.Writer.Header().Set("Origin", "*")
//
//	ctx.Next()
//}
