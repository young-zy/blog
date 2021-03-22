package models

import (
	"github.com/gin-gonic/gin/binding"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id             int    `gorm:"type:INT;NOT NULL" json:"id"`
	Username       string `gorm:"type:VARCHAR(45);NOT NULL" json:"username"`
	Email          string `gorm:"type:VARCHAR(100);NOT NULL" json:"email"`
	HashedPassword string `gorm:"type:VARCHAR(300);NOT NULL" json:"-"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type UserUpdate struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"required,password"`
	NewPassword string `json:"newPassword" binding:"omitempty,password"`
}

func init() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = validate.RegisterValidation("username", username)
		_ = validate.RegisterValidation("password", password)
	}
}

func username(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var length int
	for i, c := range username {
		if i == 0 && (unicode.IsNumber(c) || c == '_') {
			return false
		}
		length++
		switch {
		case unicode.IsNumber(c) || unicode.IsLetter(c) || c == '-' || c == '_' || c == '@':
			continue
		default:
			return false
		}
	}
	return length >= 4 && length <= 20
}

func password(fl validator.FieldLevel) bool {
	// contains at least one capitalized, one uncapitalized, one number, length at least 8 but no more than 32
	pass := fl.Field().String()
	var (
		length int
		number bool
		upper  bool
		normal bool
	)
	for _, c := range pass {
		length++
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			continue
		case unicode.IsLetter(c):
			normal = true
		default:
			return false
		}
	}
	return length > 8 && length < 32 && normal && upper && number
}
