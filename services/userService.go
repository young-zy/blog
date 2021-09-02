package services

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/databases"
	"github.com/young-zy/blog/models"
)

// Register is used for registering a new user
func Register(c *gin.Context, username, password, email string) bool {
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

// GetUser acquire user by username
func GetUser(c *gin.Context, username string) (user *models.User, ok bool) {
	ok = true
	user, err := databases.Default.GetUser(c, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		ok = false
	}
	return
}

// SetAvatar sets the avatar of a user
func SetAvatar(c *gin.Context, username, avatar string) (ok bool) {
	ok = false
	tx := databases.GetTransaction()
	user, err := tx.GetUser(c, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		tx.Rollback()
		return
	}
	user.Avatar = avatar
	err = tx.UpdateUser(c, user)
	if err != nil {
		common.NewInternalError(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	ok = true
	return
}

// UpdateUser updates the user info
func UpdateUser(c *gin.Context, user *models.User) {
	operator, exists := c.Get("User")
	if !exists {
		common.NewInternalError(c, errors.New("user object not found in context"))
		return
	}
	perm, err := enforcer.Enforce(operator.(models.User), user, "updateUser")
	if err != nil || !perm {
		_ = c.Error(common.NewForbiddenError("permission denied")).SetType(gin.ErrorTypePublic)
		return
	}
	err = databases.Default.UpdateUser(c, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user to be updated not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
	}
}

// DeleteUser deletes the user
func DeleteUser(c *gin.Context, userID *uint) {
	err := databases.Default.DeleteUser(c, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("user to be deleted not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
	}
}
