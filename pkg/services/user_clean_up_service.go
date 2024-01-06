package services

import (
	"fmt"
	"log"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/utils"
)

func DeleteUser(userId uint) error {
	tx := initializers.DB.Begin()

	// delete user charts once added too

	err := models.HandleTemplateSplitDelete(tx, userId, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.HandleTemplateSplitWorkoutLinkDelete(tx, userId, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.HandleTemplateWorkoutDelete(tx, nil, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.HandleTemplateExerciseDelete(tx, nil, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	condition := "user_id = ?"
	err = models.HandleCalendarWorkoutDelete(tx, nil, userId, condition)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.HandleCalendarExerciseDelete(tx, nil, nil, userId, condition)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.HandleCalendarSetDelete(tx, nil, nil, userId, condition)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.DeleteAllRateLimitRecordsForUser(tx, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.DeleteUser(tx, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
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
