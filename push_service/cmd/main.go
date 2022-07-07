package main

import (
	"context"
	"github.com/pavel/push_service/config"
	k "github.com/pavel/push_service/pkg/broker/kafka"
	"github.com/pavel/push_service/pkg/logger"
	"github.com/pavel/push_service/pkg/service/socket"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	log := logger.GetLogger()
	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Fatalf("Error init config: ERROR: %s", err.Error())
	}
	broadcast := make(chan socket.Broadcast)
	hub := socket.NewHub(broadcast)
	kafka := k.InitKafkaBrokerReader(cfg)

	go hub.Run()
	go kafka.Read(context.Background(), hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		socket.ServeWs(hub, w, r, r.URL.Query().Get("username"))
	})
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write([]byte("Hello from push service: your username:" + r.URL.Query().Get("username")))
	})
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Infof("Init listener adr: %s%s", cfg.Server.Host, cfg.Server.Port)
	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, nil)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}
