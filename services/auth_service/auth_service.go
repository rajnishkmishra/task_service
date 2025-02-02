package auth_service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/task_service/services/user_service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type AuthService struct {
	userService *user_service.UserService
}

func NewAuthService(userService *user_service.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (a *AuthService) AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenHeader := c.Request.Header.Get("token")
		if accessTokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}

		accessToken := strings.Split(accessTokenHeader, " ")
		if len(accessToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"message": "unauthorized access"})
			return
		}

		mapClaims := jwt.MapClaims{}
		GlobalAuthKey := []byte(viper.GetString("jwt.auth_key"))
		token, err := jwt.ParseWithClaims(accessToken[1], mapClaims, func(jwtT *jwt.Token) (interface{}, error) {
			if _, ok := jwtT.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("an error occured")
			}
			return GlobalAuthKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		if claims["uid"] == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}

		if claims["iat"] == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}

		if claims["exp"] == nil || (claims["exp"] != nil && time.Now().Unix() > int64(claims["exp"].(float64))) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}

		uid := claims["uid"].(float64)
		userId := uint64(uid)
		user, err := a.userService.GetUserDetailsByID(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}
		if user.LoginTime.Unix() < int64(claims["iat"].(float64)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "unauthorized access"})
			return
		}
		c.Next()
	}
}

func (a *AuthService) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringIP := c.Request.Header.Get("X-Forwarded-For")

		blockedIPMutex.Lock()
		_, isPresent := blockedIPs[stringIP]
		blockedIPMutex.Unlock()

		if isPresent {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, map[string]interface{}{"error": "Too many requests, please try again later."})
			return
		}

		ipMapMutex.Lock()
		defer ipMapMutex.Unlock()
		_, isPresent = ipMap[stringIP]

		if isPresent {
			ipMap[stringIP].Enqueue(time.Now())
			ipMap[stringIP].Format(time.Now())
			if ipMap[stringIP].IsThresholdReached() {
				blockedIPMutex.Lock()
				blockedIPs[stringIP] = true
				blockedIPMutex.Unlock()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, map[string]interface{}{"error": "Too many requests, please try again later."})
				return
			}
		} else {
			ipMap[stringIP] = NewRateLimiter()
			ipMap[stringIP].Enqueue(time.Now())
		}

		c.Next()
	}
}
