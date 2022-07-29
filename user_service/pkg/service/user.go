package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/pavel/user_service/pkg/utils"
	"mime/multipart"
	"net/url"
)

type User interface {
	GetUser(userId uint64) (error, *model.SelectUser)
	FindUser(filer model.UserFilter) (error, []*model.SelectUser)
	UpdateAvatar(ctx *gin.Context, userId uint64, photo multipart.File) (error, *model.SelectUser)
	UpdateUser(user model.UpdateUser) (error, *model.SelectUser)
}

type UserService struct {
	repo repository.User
	cfg  *config.Config
	log  *logger.Logger
}

func (u UserService) UpdateUser(user model.UpdateUser) (error, *model.SelectUser) {
	return u.repo.Update(user)
}

func (u UserService) UpdateAvatar(ctx *gin.Context, userId uint64, photo multipart.File) (error, *model.SelectUser) {
	awsSession := ctx.MustGet("aws_session").(*session.Session)
	aws := utils.InitAwsService(u.cfg, u.log, awsSession)
	srcImg, err := aws.UploadFileToS3(photo, fmt.Sprintf("user/%d", userId), 500, 500)
	if err != nil {
		return err, nil
	}
	err, user, oldAvatar := u.repo.UpdateAvatar(userId, srcImg)
	if err != nil {
		return err, user
	}
	if oldAvatar != "" {
		u, _ := url.Parse(oldAvatar)
		aws.DeleteFileFromS3(u.Path)
	}

	return err, user
}

func InitUserService(repo repository.User, cfg *config.Config, log *logger.Logger) UserService {
	return UserService{
		repo: repo,
		cfg:  cfg,
		log:  log,
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
