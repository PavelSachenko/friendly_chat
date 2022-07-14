package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/twinj/uuid"
	"strconv"
	"time"
)

type Auth interface {
	SigIn(user *model.User) (error, *model.TokenDetails)
	SignUp(user *model.User, roleId uint32) (error, *model.TokenDetails)
	Logout(accessToken string) error
	RefreshToken(refreshToken string) (error, *model.TokenDetails)
	CheckAuthorization(accessToken string) (error, uint64)
	GetUserIdByRefreshToken(refreshToken string) (error, uint64)
}

type AuthService struct {
	repo   repository.Auth
	config *config.Config
}

func InitAuthService(db repository.Auth, cfg *config.Config) AuthService {
	return AuthService{
		repo:   db,
		config: cfg,
	}
}

func (as AuthService) SigIn(user *model.User) (error, *model.TokenDetails) {
	err, userId := as.repo.IsUserExist(as.hashPassword(user.PasswordHash), user.Username)
	user.ID = userId
	if err != nil {
		return err, nil
	}
	return as.createToken(user.ID)
}

func (as AuthService) SignUp(user *model.User, roleId uint32) (error, *model.TokenDetails) {
	user.PasswordHash = as.hashPassword(user.PasswordHash)
	err, id := as.repo.CreateUser(user, roleId)
	if err != nil {
		return err, nil
	}
	return as.createToken(id)
}

func (as *AuthService) hashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	sha.Write([]byte(as.config.UserPasswordHashSalt))

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (as AuthService) Logout(accessToken string) error {
	tokenAuth, err := as.extractTokenMetadata(accessToken)
	if err != nil {
		return err
	}
	err = as.repo.DeleteToken(tokenAuth.AccessUuid)
	if err != nil {
		return err
	}
	return nil
}

func (as AuthService) RefreshToken(refreshToken string) (error, *model.TokenDetails) {
	err, refreshUUID, userId := as.getRefreshUUIDAndUserId(refreshToken)
	if err != nil {
		return err, nil
	}
	delErr := as.repo.DeleteToken(refreshUUID)
	if delErr != nil { //if any goes wrong
		return errors.New("unauthorized"), nil
	}

	return as.createToken(userId)
}

func (as AuthService) GetUserIdByRefreshToken(refreshToken string) (error, uint64) {
	err, _, userId := as.getRefreshUUIDAndUserId(refreshToken)
	if err != nil {
		return err, 0
	}
	return nil, userId
}

func (as AuthService) getRefreshUUIDAndUserId(refreshToken string) (error, string, uint64) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.config.Auth.AuthRefreshTokenSalt), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return errors.New("Refresh token expired"), "", 0
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err, "", 0
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return err, "", 0
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return err, "", 0
		}

		return nil, refreshUuid, userId
	} else {
		return errors.New("refresh expired"), "", 0
	}
}

func (as AuthService) CheckAuthorization(accessToken string) (error, uint64) {
	tokenAuth, err := as.extractTokenMetadata(accessToken)
	if err != nil {
		return err, 0
	}
	err, userId := as.repo.GetUserIdFromToken(tokenAuth.AccessUuid)
	if err != nil {
		return err, 0
	}
	return nil, userId
}

func (as *AuthService) createToken(userId uint64) (error, *model.TokenDetails) {
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(as.config.Auth.AuthAccessTokenExpire).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(as.config.Auth.AuthRefreshTokenExpire).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(as.config.Auth.AuthAccessTokenSalt))
	if err != nil {
		return err, nil
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(as.config.Auth.AuthRefreshTokenSalt))
	if err != nil {
		return err, nil
	}

	err = as.repo.SaveToken(userId, td)
	if err != nil {
		return err, nil
	}

	return nil, td
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func (as *AuthService) extractTokenMetadata(accessToken string) (*AccessDetails, error) {
	token, err := as.verifyToken(accessToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("Bad jwt token")
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func (as *AuthService) verifyToken(accessToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.config.Auth.AuthAccessTokenSalt), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
