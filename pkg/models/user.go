package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
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

func EncryptPassword(password string, userId uint) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil && userId == 0 {
		log.Print(fmt.Sprintf("Password hash failed for new user: %s", err.Error()))
		return "", err
	}
	if err != nil && userId != 0 {
		log.Print(fmt.Sprintf("Password hash failed for user with ID: %d: %s", userId, err.Error()))
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
	if result.Error != nil {
		log.Print(fmt.Sprintf("User creation failed for new email: %s: %s", user.Email, result.Error.Error()))
		return user, result.Error
	}

	return user, nil
}
