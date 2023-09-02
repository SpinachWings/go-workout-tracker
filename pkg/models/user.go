package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"workout-tracker-go-app/pkg/initializers"
)

type User struct {
	gorm.Model
	Email                      string `gorm:"unique;not null"`
	Password                   string `gorm:"not null"`
	EmailVerificationCode      string
	IsVerified                 bool `gorm:"not null;default:false"`
	PasswordResetCode          string
	PasswordResetCodeCreatedAt time.Time
}

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateUser(email string, hash string, verificationCode string) (User, error) {
	user := User{
		Email:                      email,
		Password:                   hash,
		EmailVerificationCode:      verificationCode,
		IsVerified:                 false,
		PasswordResetCode:          "",
		PasswordResetCodeCreatedAt: time.Now(),
	}

	result := initializers.DB.Create(&user)
	return user, result.Error
}
