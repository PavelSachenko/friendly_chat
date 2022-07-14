package main

import (
	"context"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pavel/push_service/config"
	k "github.com/pavel/push_service/pkg/broker/kafka"
	"github.com/pavel/push_service/pkg/logger"
	"github.com/pavel/push_service/pkg/service/socket"
	"io/ioutil"
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
	test := gin.New()

	test.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1000", "http://localhost:10000", "http://localhost:10001"},
		//AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers", "*"},
		//AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowWildcard:    false,
		AllowWebSockets:  true,

		MaxAge: 84000,
	}))

	test.GET("/push", func(ctx *gin.Context) {
		//cookie, err := ctx.Cookie("refresh_token")
		//if err != nil {
		//	log.Fatalf("ERROR connect to socket: %v", err)
		//	ctx.JSON(422, `{"error"": "refresh_token cookie doesn't found'"}`)
		//	return
		//}
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:10000/api/user", nil)
		if err != nil {
			ctx.JSON(422, err.Error())
		}
		req.Header.Set("Authorization", ctx.Request.Header.Get("Authorization"))
		res, err := client.Do(req)
		log.Println(req.Header)
		if err != nil {
			ctx.JSON(422, err.Error())
		}
		test, _ := ioutil.ReadAll(res.Body)
		var e struct {
			ID uint64 `json:"id"`
		}
		json.Unmarshal(test, &e)
		ctx.JSON(res.StatusCode, e.ID)
	})

	test.GET("/ws", func(ctx *gin.Context) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:10000/api/user", nil)
		if err != nil {
			ctx.JSON(422, err.Error())
		}
		param, _ := ctx.GetQuery("accessToken")

		req.Header.Set("Authorization", "Bearer "+param)
		res, err := client.Do(req)
		log.Println(req.Header)
		if err != nil {
			ctx.JSON(422, err.Error())
		}
		test, _ := ioutil.ReadAll(res.Body)
		var e struct {
			ID uint64 `json:"id"`
		}
		json.Unmarshal(test, &e)
		log.Println(e.ID)
		log.Println("Bearer " + param)

		socket.ServeWs(hub, ctx.Writer, ctx.Request, e.ID)
	})
	log.Infof("Init listener adr: %s%s", cfg.Server.Host, cfg.Server.Port)
	err = http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, test)
	if err != nil {
		log.Fatalf("Error init net listener: ERROR: %s", err.Error())
	}
}
