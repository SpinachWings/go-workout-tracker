package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
)

type sendPasswordResetEmailBody struct {
	Email string `json:"email" binding:"required"`
}

type passwordResetBody struct {
	Email            string `json:"email" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
}

func SendPasswordResetEmail(c *gin.Context) {
	var body sendPasswordResetEmailBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding user with email: %s: %s", body.Email, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "password reset email has been sent if user with this email exists"})
		return
	}

	verificationCode := randstr.String(40)

	user.PasswordResetCode = verificationCode
	user.PasswordResetCodeCreatedAt = time.Now()
	result = initializers.DB.Save(&user)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to save user with password reset code for user with ID: %d: %s", user.ID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	err = services.SendResetPasswordEmail(verificationCode, user.Email)
	if err != nil {
		log.Print(fmt.Sprintf("Error sending password reset email for user with ID: %d: %s", user.ID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	models.CreateAudit(constants.GetAuditTypes().SendPasswordResetEmail, user.ID, "")

	c.JSON(http.StatusOK, gin.H{"message": "password reset email has been sent if user with this email exists"})
}

func ResetPassword(c *gin.Context) {
	var body passwordResetBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding user with email: %s: %s", body.Email, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if user.ID == 0 || user.PasswordResetCode != body.VerificationCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code or email"})
		return
	}

	hash, err := models.EncryptPassword(body.NewPassword, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	user.Password = hash
	user.PasswordResetCode = ""
	result = initializers.DB.Save(user)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to save user with ID: %d with new password: %s", user.ID, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	models.CreateAudit(constants.GetAuditTypes().ResetPassword, user.ID, "")

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}
