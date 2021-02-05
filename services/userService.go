package services

import (
	"blog/common"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"blog/databases"
	"blog/models"
)

func Register(username, password, email string) common.HttpError {
	// TODO validate field

	// generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return common.NewBadRequestError("error creating hashed password")
	}
	user := &models.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
	}
	// insert user into database
	if err := databases.AddUser(user); err != nil {
		// check if error is mysql error
		mySQLError, ok := err.(*mysql.MySQLError)
		if ok {
			// duplicate entry error
			if mySQLError.Number == 1062 {
				return common.NewSelfDefinedError(http.StatusConflict, "username or email already exists")
			}
		}
		return common.NewInternalServerError(err.Error())
	}
	// return no content if success
	return nil
}

// acquire user by username
func GetUser(username string) (user *models.User, httpError common.HttpError) {
	user, err := databases.GetUser(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("user not found")
		} else {
			httpError = common.NewInternalServerError(err.Error())
		}
		return
	}
	return
}

func UpdateUser(user *models.User) (httpError common.HttpError) {
	err := databases.UpdateUser(user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("user to be updated not found")
		} else {
			httpError = common.NewInternalServerError(err.Error())
		}
		return
	}
	return nil
}

func DeleteUser(userId int) (httpError common.HttpError) {
	err := databases.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("user to be deleted not found")
		} else {
			httpError = common.NewInternalServerError(err.Error())
		}
		return
	}
	return
}
