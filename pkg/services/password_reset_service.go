package services

import (
	"fmt"
	"os"
	"workout-tracker-go-app/pkg/utils"
)

func SendResetPasswordEmail(verificationCode string, userEmail string) error {
	// url will actually be something defined on the client eventually...
	url := os.Getenv("CLIENT_ORIGIN") + "/password/reset/confirm/" + verificationCode
	subject := "Go Workout Tracker - Password Reset"
	body := fmt.Sprintf("Click link to reset password: %s", url)
	return utils.SendEmail(subject, body, userEmail)
}
