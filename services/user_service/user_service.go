package user_service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/task_service/models"
	"bitbucket.org/task_service/models/vm"
	"bitbucket.org/task_service/utils"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) Login(ctx *utils.Context, request vm.LoginRequest) (response vm.LoginResponse, werr utils.WrapperError) {
	if strings.TrimSpace(request.PhoneNumber) == "" {
		werr = utils.NewWrapperError(http.StatusBadRequest, errors.New("invalid phone number"))
		return
	}

	var user models.User
	err := u.db.Table("users").Where("phone_number = ?", request.PhoneNumber).First(&user).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		if err == gorm.ErrRecordNotFound {
			user, werr = u.MakeNewUser(ctx, request.PhoneNumber)
			if werr != nil {
				logrus.WithContext(ctx.Ctx).Error(err)
				return
			}
		} else {
			werr = utils.NewWrapperError(http.StatusInternalServerError, err)
			return
		}
	}

	accessToken, err := u.GetAccessToken(ctx, user)
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	user.LoginTime = time.Now()
	err = u.db.Table("users").Save(&user).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.UserID = user.ID
	response.AccessToken = accessToken
	return
}

func (u *UserService) MakeNewUser(ctx *utils.Context, phoneNumber string) (user models.User, werr utils.WrapperError) {
	user.PhoneNumber = phoneNumber
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := u.db.Table("users").Save(&user).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	return
}

func (u *UserService) GetAccessToken(ctx *utils.Context, user models.User) (string, error) {
	now := time.Now()
	claims := &jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(48 * time.Hour).Unix(),
		"uid": user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	GlobalAuthKey := []byte(viper.GetString("jwt.auth_key"))
	signedToken, err := token.SignedString(GlobalAuthKey)
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		return signedToken, err
	}
	return signedToken, err
}

func (u *UserService) GetUserDetailsByID(userId uint64) (user models.User, err error) {

	err = u.db.Table("users").Where("id = ?", userId).First(&user).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}
