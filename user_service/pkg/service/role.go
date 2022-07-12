package service

import (
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
)

type Role interface {
	All() (error, []*model.Role)
}

type RoleService struct {
	repo repository.Role
}

func InitRoleService(repo repository.Role) RoleService {
	return RoleService{
		repo: repo,
	}
}

func (r RoleService) All() (error, []*model.Role) {
	return r.repo.GetRoles()
}
