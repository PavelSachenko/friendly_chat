package service

import (
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
)

type User interface {
	GetUser(userId uint64) (error, *model.User)
}

type UserService struct {
	repo repository.User
}

func InitUserService(repo repository.User) UserService {
	return UserService{
		repo: repo,
	}
}

func (u UserService) GetUser(userId uint64) (error, *model.User) {
	return u.repo.One(userId)
}
