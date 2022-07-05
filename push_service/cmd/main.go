package main

import (
	"github.com/pavel/push_service/config"
	"github.com/pavel/push_service/pkg/logger"
	"github.com/pavel/push_service/pkg/service/socket"
	"net/http"
)

//
//func serveHome(w http.ResponseWriter, r *http.Request) {
//
//	log.Println(r.URL)
//	if r.URL.Path != "/" {
//		http.Error(w, "Not found", http.StatusNotFound)
//		return
//	}
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	http.ServeFile(w, r, "home.html")
//}

func main() {
	log := logger.GetLogger()

	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Fatalf("Error init config: ERROR: %s", err.Error())
	}
	hub := socket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		socket.ServeWs(hub, w, r)
	})
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("Hello from push service"))
	})
	log.Infof("Init listener adr: %s%s", cfg.Server.Host, cfg.Server.Port)
	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, nil)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}
