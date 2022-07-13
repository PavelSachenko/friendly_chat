package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/utils"
	"github.com/pavel/user_service/pkg/validation"
	"net/http"
	"time"
)

type signUpRequest struct {
	Password       string `json:"password"  validate:"required,min=5,max=500"`
	Username       string `json:"username"   validate:"required,min=2,max=500"`
	RepeatPassword string `json:"repeat_password" validate:"required,min=6,max=60,eqfield=Password"`
}

type signInRequest struct {
	Username string `json:"username"  validate:"required,min=2,max=500"`
	Password string `json:"password"  validate:"required,min=5,max=500"`
}

type SignInResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (h *Handler) signUp(ctx *gin.Context) {
	var signUpRequest signUpRequest
	requestErrors := h.validator.ValidateRequest(ctx, &signUpRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}

	err, tokens := h.auth.SignUp(&model.User{Username: signUpRequest.Username, PasswordHash: signUpRequest.Password}, 2)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, validation.IError{Value: err.Error()})
		return
	}
	setRefreshTokenInCookie(ctx, tokens)
	ctx.JSON(http.StatusCreated, &SignInResponse{RefreshToken: tokens.RefreshToken, AccessToken: tokens.AccessToken})
}

func (h *Handler) signIn(ctx *gin.Context) {
	var signInRequest signInRequest
	requestErrors := h.validator.ValidateRequest(ctx, &signInRequest)
	if requestErrors != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, requestErrors)
		return
	}
	err, tokens := h.auth.SigIn(&model.User{Username: signInRequest.Username, PasswordHash: signInRequest.Password})
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, validation.IError{Value: err.Error()})
		return
	}
	setRefreshTokenInCookie(ctx, tokens)
	ctx.JSON(http.StatusCreated, &SignInResponse{RefreshToken: tokens.RefreshToken, AccessToken: tokens.AccessToken})
}

func (h *Handler) refreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Request.Cookie("refresh_token")
	if refreshToken.Value == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, validation.IError{Field: "refresh_token", Tag: "not required", Value: "refresh token doesn't exist in cookie"})
		return
	}
	err, tokens := h.auth.RefreshToken(refreshToken.Value)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, validation.IError{Value: err.Error()})
		return
	}
	setRefreshTokenInCookie(ctx, tokens)
	ctx.JSON(http.StatusCreated, &SignInResponse{RefreshToken: tokens.RefreshToken, AccessToken: tokens.AccessToken})
}

func (h *Handler) logout(ctx *gin.Context) {
	token := utils.GetBearerToken(ctx)
	err := h.auth.Logout(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	ctx.JSON(http.StatusNoContent, nil)
}

func setRefreshTokenInCookie(ctx *gin.Context, tokens *model.TokenDetails) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 36),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})
}
