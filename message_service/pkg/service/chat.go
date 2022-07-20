package service

type Chat interface {
	AddToUsers(title, description string, userIds []int)
}

type ChatService struct {
}

func (c ChatService) AddToUsers(title, description string, userIds []int) {

}
