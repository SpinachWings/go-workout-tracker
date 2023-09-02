package services

import (
	"fmt"
	"os"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/utils"
)

func SendVerificationEmail(verificationCode string, userEmail string) error {
	//eventually url will be something defined by client...
	url := os.Getenv("CLIENT_ORIGIN") + "/verify/email/" + verificationCode
	subject := "Go Workout Tracker - Email Verification"
	body := fmt.Sprintf("Click link to verify email: %s", url)
	return utils.SendEmail(subject, body, userEmail)
}

func VerifyEmail(verificationCode string, userEmail string) error {
	var updatedUser models.User
	result := initializers.DB.First(&updatedUser, "email_verification_code = ? AND email = ?", verificationCode, userEmail)
	if result.Error != nil {
		return result.Error
	}

	updatedUser.EmailVerificationCode = ""
	updatedUser.IsVerified = true
	result = initializers.DB.Save(&updatedUser)
	return result.Error
}
