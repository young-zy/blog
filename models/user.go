package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id             int    `gorm:"type:INT;NOT NULL" json:"id"`
	Username       string `gorm:"type:VARCHAR(45);NOT NULL" json:"username"`
	Email          string `gorm:"type:VARCHAR(100);NOT NULL" json:"email"`
	HashedPassword string `gorm:"type:VARCHAR(300);NOT NULL" json:"-"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,usernameRegex"`
	Password string `json:"password" binding:"required,passwordRegex"`
}

type UserRegister struct {
	Username string `json:"username" binding:"required,usernameRegex"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,passwordRegex"`
}

type UserUpdate struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"required,passwordRegex"`
	NewPassword string `json:"newPassword" binding:"omitempty,passwordRegex"`
}

func init() {
	validate := validator.New()
	_ = validate.RegisterValidation("usernameRegex", usernameRegex)
	_ = validate.RegisterValidation("passwordRegex", passwordRegex)
}

func usernameRegex(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile("^(?=.{4,20}$)(?![_0-9])(?!.*[_]{2})[a-zA-Z0-9_-]+(?<![_])$")
	username := fl.Field().String()
	if reg.MatchString(username) {
		return true
	}
	return false
}

func passwordRegex(fl validator.FieldLevel) bool {
	// contains at least one capitalized, one uncapitalized, one number, length at least 8 but no more than 32
	reg := regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)[a-zA-Z!@#$%^&*\\d]{8,32}$")
	pass := fl.Field().String()
	if reg.MatchString(pass) {
		return true
	}
	return false
}
