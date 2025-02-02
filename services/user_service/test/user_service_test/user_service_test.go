package user_service_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"bitbucket.org/task_service/models"
	"bitbucket.org/task_service/models/vm"
	"bitbucket.org/task_service/services/user_service"
	"bitbucket.org/task_service/utils"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.User{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	userService := user_service.NewUserService(db)

	phoneNumber := "1234567890"
	user := models.User{
		PhoneNumber: phoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&user)

	loginRequest := vm.LoginRequest{
		PhoneNumber: phoneNumber,
	}

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	response, werr := userService.Login(ctx, loginRequest)

	assert.Nil(t, werr)

	assert.Equal(t, user.ID, response.UserID)

	authKey := []byte(viper.GetString("jwt.auth_key"))
	_, err := jwt.Parse(response.AccessToken, func(jwtT *jwt.Token) (interface{}, error) {
		if _, ok := jwtT.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("an error occured")
		}
		return authKey, nil
	})
	assert.Nil(t, err)

	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	assert.NotNil(t, updatedUser.LoginTime)
}

func TestMakeNewUser(t *testing.T) {
	userService := user_service.NewUserService(db)

	phoneNumber := "0987654321"

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	user, werr := userService.MakeNewUser(ctx, phoneNumber)

	assert.Nil(t, werr)
	assert.NotNil(t, user.ID)
	assert.Equal(t, phoneNumber, user.PhoneNumber)

	var dbUser models.User
	db.First(&dbUser, user.ID)
	assert.Equal(t, phoneNumber, dbUser.PhoneNumber)
}

func TestGetAccessToken(t *testing.T) {
	userService := user_service.NewUserService(db)

	user := models.User{
		PhoneNumber: "1234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&user)

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	accessToken, err := userService.GetAccessToken(ctx, user)

	assert.Nil(t, err)

	authKey := []byte(viper.GetString("jwt.auth_key"))
	_, parseErr := jwt.Parse(accessToken, func(jwtT *jwt.Token) (interface{}, error) {
		if _, ok := jwtT.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("an error occured")
		}
		return authKey, nil
	})
	assert.Nil(t, parseErr)
}

func TestGetUserDetailsByID(t *testing.T) {
	userService := user_service.NewUserService(db)

	user := models.User{
		PhoneNumber: "1234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&user)

	resultUser, err := userService.GetUserDetailsByID(user.ID)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, resultUser.ID)
	assert.Equal(t, user.PhoneNumber, resultUser.PhoneNumber)
}
