package services

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"blog/common"
	"blog/databases"
	"blog/models"
)

func Register(c *gin.Context, username, password, email string) bool {
	// TODO validate field
	// generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		_ = c.Error(common.NewBadRequestError("error creating hashed password")).SetType(gin.ErrorTypePublic)
	}
	user := &models.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
	}
	// insert user into database
	if err := databases.Default.AddUser(c, user); err != nil {
		// check if error is mysql error
		mySQLError, ok := err.(*mysql.MySQLError)
		if ok {
			// duplicate entry error
			if mySQLError.Number == 1062 {
				_ = c.Error(common.NewSelfDefinedError(http.StatusConflict, "username or email already exists")).
					SetType(gin.ErrorTypePublic)
			}
		} else {
			common.NewInternalError(c, err)
		}
		return false
	}
	return true
}

// acquire user by username
func GetUser(c *gin.Context, username string) (user *models.User, ok bool) {
	ok = true
	user, err := databases.Default.GetUser(c, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user not found"))
		} else {
			common.NewInternalError(c, err)
		}
		ok = false
	}
	return
}

func UpdateUser(c *gin.Context, user *models.User) {
	err := databases.UpdateUser(user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user to be updated not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
	}
}

func DeleteUser(c *gin.Context, userId int) {
	err := databases.Default.DeleteUser(c, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user to be deleted not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
	}
}
