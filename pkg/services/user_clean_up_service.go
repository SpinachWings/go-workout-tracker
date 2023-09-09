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
	// incomplete fn
	// ensure all correlated tables are deleted as well
	var user models.User
	initializers.DB.First(&user, userId)
	initializers.DB.Unscoped().Delete(&user)
}

func DeleteExpiredUnverifiedUsers() {
	timeFrom := utils.CurrentTimeMinusHoursAsTime(constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.ExpiryTimeInHours)
	sleepTime := constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.SleepTimeInHours

	var usersWithExpiredUnverifiedEmails []models.User
	result := initializers.DB.Find(&usersWithExpiredUnverifiedEmails, "is_verified = ? AND created_at <= ?", false, timeFrom)

	if len(usersWithExpiredUnverifiedEmails) == 0 {
		log.Print("Ran delete expired unverified users but no unverified users were found")
		utils.SleepForHours(sleepTime)
		DeleteExpiredUnverifiedUsers()
	}

	if result.Error != nil {
		log.Print(fmt.Sprintf("Error finding unverified users to delete: %s"), result.Error.Error())
		utils.SleepForHours(sleepTime)
		DeleteExpiredUnverifiedUsers()
	}

	result = initializers.DB.Unscoped().Delete(&usersWithExpiredUnverifiedEmails)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Error deleting unverified users: %s"), result.Error.Error())
	}

	log.Print("Successfully deleted unverified users unverified users")
	utils.SleepForHours(sleepTime)
	DeleteExpiredUnverifiedUsers()
}

func RemoveExpiredPasswordResetCodes() {
	timeFrom := utils.CurrentTimeMinusHoursAsTime(constants.GetExpiryCheckTimes().UserPasswordResetCode.ExpiryTimeInHours)
	sleepTime := constants.GetExpiryCheckTimes().UserPasswordResetCode.SleepTimeInHours

	var usersWithPasswordResetCodes []models.User
	result := initializers.DB.Find(&usersWithPasswordResetCodes, "password_reset_code != ? AND password_reset_code_created_at <= ?", "", timeFrom)

	if len(usersWithPasswordResetCodes) == 0 {
		log.Print("Ran remove expired password reset codes but no users with expired password reset codes were found")
		utils.SleepForHours(sleepTime)
		RemoveExpiredPasswordResetCodes()
	}

	if result.Error != nil {
		log.Print(fmt.Sprintf("Error finding users with expired password reset codes to delete: %s"), result.Error.Error())
		utils.SleepForHours(sleepTime)
		RemoveExpiredPasswordResetCodes()
	}

	for _, user := range usersWithPasswordResetCodes {
		user.PasswordResetCode = ""
	}

	result = initializers.DB.Save(&usersWithPasswordResetCodes)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to remove users' expired password reset codes: %s", result.Error.Error()))
	}

	log.Print("Successfully removed users' expired password reset codes")
	utils.SleepForHours(sleepTime)
	RemoveExpiredPasswordResetCodes()
}
