package services

import (
	"fmt"
	"log"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/utils"
)

func DeleteUser(userId uint) {
	// ensure all correlated tables are deleted as well!

	var user models.User
	// handle error
	initializers.DB.First(&user, userId)
	initializers.DB.Delete(&user)
}

func InitCheckForExpiredUnverifiedUsers() {
	go DeleteExpiredUnverifiedUsers()
}

func DeleteExpiredUnverifiedUsers() {
	// we don't log some of the full errors as it will error if no users that meet the condition are found - which is expected sometimes

	var usersWithExpiredUnverifiedEmails []models.User
	result := initializers.DB.Delete(&usersWithExpiredUnverifiedEmails, "is_verified = ? AND created_at <= ?", false, constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.ExpiryTimeInHours)
	if result.Error != nil {
		log.Print("Failed to delete expired users with unverified emails, or none exist")
	} else {
		log.Print("Deleted expired users with unverified emails")
	}

	utils.SleepForHours(constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.SleepTimeInHours)
	DeleteExpiredUnverifiedUsers()
}

func InitCheckForExpiredPasswordResetCodes() {
	go RemoveExpiredPasswordResetCodes()
}

func RemoveExpiredPasswordResetCodes() {
	// we don't log some of the full errors as it will error if no users that meet the condition are found - which is expected sometimes

	var usersWithPasswordResetCodes []models.User
	result := initializers.DB.Find(&usersWithPasswordResetCodes, "password_reset_code != ?", "")
	if result.Error != nil {
		log.Print("Failed to find users with password reset codes, or none exist")
	} else {
		containsExpired := false
		for _, userWithPasswordResetCode := range usersWithPasswordResetCodes {
			passwordResetCodeIsExpired := userWithPasswordResetCode.PasswordResetCodeCreatedAt.Before(utils.CurrentTimeMinusHoursAsTime(constants.GetExpiryCheckTimes().UserPasswordResetCode.ExpiryTimeInHours))
			if passwordResetCodeIsExpired {
				userWithPasswordResetCode.PasswordResetCode = ""
				containsExpired = true // ensure this actually gets re-assigned?
			}
		}

		if containsExpired {
			result := initializers.DB.Save(&usersWithPasswordResetCodes)
			if result.Error != nil {
				log.Print(fmt.Sprintf("Failed to remove users' expired password reset codes: %s", result.Error.Error()))
			}
			log.Print("Removed users' expired password reset codes")
		}
	}

	utils.SleepForHours(constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.SleepTimeInHours)
	RemoveExpiredPasswordResetCodes()
}
