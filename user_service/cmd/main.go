package main

import (
	"github.com/pavel/user_service/config"
	api "github.com/pavel/user_service/pkg/api/http"
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
	//http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Write([]byte("Hello from user service"))
	//})
	//err, postgres := db.InitPostgres(cfg)
	//if err != nil {
	//	log.Fatalf("Error init postgres db: CONFIG: %v ERROR: %s", cfg.DB, err.Error())
	//}
	//err, redis := db.InitRedis(cfg)
	//if err != nil {
	//	log.Fatalf("Error init postgres db: CONFIG: %v ERROR: %s", cfg.DB, err.Error())
	//}
	userRepo := repository.InitUserPostgres(nil, log)
	authRepo := repository.InitAuthRedis(nil, nil, log)

	userService := service.InitUserService(userRepo)
	authService := service.InitAuthService(authRepo, cfg)

	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, api.InitHandler(*log, userService, authService).Handle())
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}
