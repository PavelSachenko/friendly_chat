package service

import (
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/repository"
)

type Message interface {
	Send(senderUserId, chatId int, body string) (int, error)
	GetAllForChat(filterMessage model.FilterMessage) ([]*model.MessageChat, error)
}
type MessageService struct {
	repo repository.Message
}

func InitMessageService(repo repository.Message) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (m MessageService) Send(senderUserId, chatId int, body string) (int, error) {
	return m.repo.CreateChatMessage(senderUserId, chatId, body)
}

func (m MessageService) GetAllForChat(filterMessage model.FilterMessage) ([]*model.MessageChat, error) {
	return m.repo.All(filterMessage)
}
