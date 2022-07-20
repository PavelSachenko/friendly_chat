package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pavel/message_service/config"
	"github.com/pavel/message_service/gateway/point/command"
	"github.com/pavel/message_service/pkg/db"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/repository"
)

func main() {
	log := logger.GetLogger()
	log.Infof("Init config")
	err, cfg := config.InitConfig(log)
	if err != nil {
		log.Fatalf("Error init config: ERROR: %s", err.Error())
	}
	err, postgre := db.InitPostgres(cfg, db.InitPostgresQueryBuilder())
	if err != nil {
		log.Fatalf("Error init postgre: ERROR: %s", err.Error())
	}

	chat := repository.InitChatPostgreSQL(*postgre, *log)
	test, err := chat.All(model.FilterChat{UserId: 6, Limit: 20, Offset: 0})
	//fmt.Printf("%v\n%v\n", test[0], test[0].LastMessage)
	fmt.Printf("%v\n", test)
	//fmt.Printf("%v\n%v\n", test[1], test[1].LastMessage)
	//titles := []model.UserTitleChats{{UserId: 6, Title: "User 10"}, {UserId: 10, Title: "User 6"}}
	//_ = chat.CreatePrivate("", titles)
	//owners := []model.UsersOwners{{UserId: 1, IsOwner: true}, {UserId: 2, IsOwner: false}, {UserId: 3, IsOwner: false}}
	//chat.CreateGroup("Group chat", "test", owners)

	//message := repository.InitMessagePostgeSQL(postgre, log)
	//message.CreateMessage(6, 3, "Hello User 10")
	return
	r := gin.New()
	api := r.Group("/api/message_service")
	command.RegisterHandlers(api, cfg)

	log.Fatal(r.Run(cfg.Server.GatewayAddress))
}
