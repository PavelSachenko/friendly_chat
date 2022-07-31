package service

import (
	"encoding/json"
	"fmt"
	"github.com/pavel/message_service/pkg/broker"
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/repository"
	"golang.org/x/net/context"
	"strings"
)

type Message interface {
	Send(senderUserId, chatId int, body string) (int, error)
	GetAllForChat(filterMessage model.FilterMessage) ([]*model.MessageChat, error)
}
type MessageService struct {
	messageRepo  repository.Message
	chatRepo     repository.Chat
	brokerWriter broker.BrokerWriter
}

func InitMessageService(messageRepo repository.Message, chatRepo repository.Chat, brokerWriter broker.Broker) *MessageService {
	return &MessageService{
		messageRepo:  messageRepo,
		chatRepo:     chatRepo,
		brokerWriter: brokerWriter,
	}
}

func (m MessageService) Send(senderUserId, chatId int, body string) (int, error) {
	message, err := m.messageRepo.CreateChatMessage(senderUserId, chatId, body)
	if err != nil {
		return 0, err
	}
	userIds, err := m.chatRepo.GetChatUserIds(chatId)
	if err != nil {
		return 0, err
	}
	strUserIds := strings.Trim(strings.Replace(fmt.Sprint(userIds), " ", ",", -1), "[]")
	messageJson, err := json.Marshal(message)
	if err != nil {
		return 0, err
	}
	go m.brokerWriter.Push(context.Background(), nil, []byte(
		fmt.Sprintf(`{"user_ids": [%s], "type": "sms", "data": %s}`, strUserIds, messageJson),
	),
	)

	return message.ID, err
}

func (m MessageService) GetAllForChat(filterMessage model.FilterMessage) ([]*model.MessageChat, error) {
	return m.messageRepo.All(filterMessage)
}
