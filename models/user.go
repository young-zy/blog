package models

import (
	"database/sql/driver"
	"errors"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Role role for user
type Role int

const (
	// enumeration of roles
	RoleUser Role = iota
	RoleAdmin
)

var roleMap = map[string]Role{
	"user":  RoleUser,
	"admin": RoleAdmin,
}

// Scan is a scanner for Role when mapping from database
func (r *Role) Scan(value interface{}) error {
	str, ok := value.([]uint8)
	if !ok {
		return errors.New("failed to parse value to string")
	}
	*r, ok = roleMap[string(str)]
	if !ok {
		return errors.New("unknown role received")
	}
	return nil
}

// Value returns the value for Role when saving to database
func (r Role) Value() (driver.Value, error) {
	return r.String(), nil
}

// String is used to map Role to string
func (r *Role) String() string {
	return [...]string{"user", "admin"}[*r]
}

// User is used for base and orm
type User struct {
	ID             *uint  `gorm:"type:INT;NOT NULL" json:"id"`
	Username       string `gorm:"type:VARCHAR(45);NOT NULL" json:"username"`
	Email          string `gorm:"type:VARCHAR(100);NOT NULL" json:"email"`
	HashedPassword string `gorm:"type:VARCHAR(300);NOT NULL" json:"-"`
	Role           Role   `gorm:"type:VARCHAR(45);NOT NULL" json:"role"`
	Avatar         string `gorm:"type:MEDIUMTEXT" json:"avatar"`
}

// SimpleUser is used tobe embedded in other structs
type SimpleUser struct {
	ID       *uint  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// GetSimpleUser is used to map User to a SimpleUser
func (u *User) GetSimpleUser() *SimpleUser {
	return &SimpleUser{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

// LoginRequest is used when log in
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegister is used when register
type UserRegister struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

// UserUpdate is used when update
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
