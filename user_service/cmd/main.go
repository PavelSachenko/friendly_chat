package main

import (
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/logger"
	"net/http"
)

func main() {
	log := logger.GetLogger()

	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Fatalf("Error init config: ERROR: %s", err.Error())
	}
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("Hello from user service"))
	})

	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, nil)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}
