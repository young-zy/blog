package models

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
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
