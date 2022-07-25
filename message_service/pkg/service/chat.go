package service

import (
	"github.com/pavel/message_service/pkg/model"
	"github.com/pavel/message_service/pkg/repository"
)

type Chat interface {
	CreateBetweenUsers(chats []model.UserTitleChats) (int, error)
	CreateGroup(description, title string, UsersOwners []model.UsersOwners) (int, error)
	GetAllForUser(filterChat model.FilterChat) ([]*model.AllChat, error)
}
type ChatService struct {
	repo repository.Chat
}

func InitChatService(chat repository.Chat) *ChatService {
	return &ChatService{
		repo: chat,
	}
}

func (c ChatService) CreateBetweenUsers(chats []model.UserTitleChats) (int, error) {
	return c.repo.CreatePrivate(chats)
}

func (c ChatService) CreateGroup(description, title string, UsersOwners []model.UsersOwners) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChatService) GetAllForUser(filterChat model.FilterChat) ([]*model.AllChat, error) {
	return c.repo.All(filterChat)
}
