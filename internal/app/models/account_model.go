package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Phone     string
	Password  string
	Role      string
	Schedules []Schedule
}

// LoginRequest : Account login struct
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type AccountResponse struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Role  string `json:"role,omitempty"`
}

// JWT : JWT token struct
type JWT struct {
	ID          uint   `json:"id,omitempty"`
	AccessToken string `json:"access_token"`
}

// Claims : Details required to identify an Admin
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type User struct {
	ID   uint
	Name string
	Role string
}

func (User) TableName() string {
	return "accounts"
}

type NameRole struct {
	Name string
	Role string
}

type EditUserRequest struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty" binding:"e164"`
	Role  string `json:"role,omitempty"`
}

type UserIDNameRoleMap map[uint]NameRole
