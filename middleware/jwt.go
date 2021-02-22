package middleware

import (
	"errors"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"blog/conf"
	"blog/models"
	"blog/services"
)

var AuthMiddleware *jwt.GinJWTMiddleware

func init() {

	config := conf.Config

	var err error
	// the jwt middleware
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "",
		Key:             []byte(config.Server.JwtKey),
		Timeout:         time.Hour * 24,
		MaxRefresh:      time.Hour * 168,
		IdentityKey:     "User",
		PayloadFunc:     payload,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizator:    authorizer,
		Unauthorized:    unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"Id":       v.Id,
			"Username": v.Username,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.User{
		Id:       int(claims["Id"].(float64)),
		Username: claims["Username"].(string),
	}
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginRequest models.LoginRequest
	// check
	if err := c.ShouldBind(&loginRequest); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	// retrieve user from database and check password
	if user, ok := services.GetUser(c, loginRequest.Username); ok {
		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginRequest.Password)); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return &models.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		}, nil
	} else {
		return nil, errors.New("username does not exist")
	}
}

func authorizer(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*models.User); ok && v.Username != "" {
		// reserved for future role or auth settings in context
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
