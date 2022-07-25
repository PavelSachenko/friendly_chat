package service

import (
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
)

type User interface {
	GetUser(userId uint64) (error, *model.SelectUser)
	FindUser(filer model.UserFilter) (error, []*model.SelectUser)
}

type UserService struct {
	repo repository.User
}

func InitUserService(repo repository.User) UserService {
	return UserService{
		repo: repo,
	}
}

func (u UserService) GetUser(userId uint64) (error, *model.SelectUser) {
	return u.repo.One(userId)
}

func (u UserService) FindUser(filer model.UserFilter) (error, []*model.SelectUser) {
	if filer.Limit == 0 {
		filer.Limit = 20
	}
	return u.repo.All(filer)
}
