package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/utils"
)

type RateLimitAction struct {
	gorm.Model
	UserId                uint
	AlternativeIdentifier string
	RateLimitAction       string `gorm:"not null"`
}

func CreateRateLimitRecord(rateLimitActionType constants.RateLimitActionType, userId uint, alternativeIdentifier string) {
	rateLimitAction := RateLimitAction{
		UserId:                userId,
		AlternativeIdentifier: alternativeIdentifier,
		RateLimitAction:       rateLimitActionType.Action,
	}

	result := initializers.DB.Create(&rateLimitAction)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to create rate limit record of type: %s with identifier: %d / %s: %s", rateLimitActionType.Action, userId, alternativeIdentifier, result.Error.Error()))
	}
}

func RateLimitIsExceeded(rateLimitActionType constants.RateLimitActionType, userId uint, alternativeIdentifier string) bool {
	timeFrom := utils.CurrentTimeMinusMinutesAsTime(rateLimitActionType.DurationInMinutes)

	var rateLimitActions []RateLimitAction
	if userId != 0 {
		rateLimitActions = GetRateLimitRecordsViaUserId("user_id = ? AND rate_limit_action = ? AND created_at >= ?", userId, rateLimitActionType.Action, timeFrom)
	} else {
		rateLimitActions = GetRateLimitRecordsViaAlternativeIdentifier("alternative_identifier = ? AND rate_limit_action = ? AND created_at >= ?", alternativeIdentifier, rateLimitActionType.Action, timeFrom)
	}

	ClearOldRateLimitRecords(rateLimitActionType, userId, alternativeIdentifier)

	if rateLimitActions == nil {
		return false
	}
	return len(rateLimitActions) > rateLimitActionType.RateLimit
}

func GetRateLimitRecordsViaUserId(conditions string, userId uint, action string, timeFrom time.Time) []RateLimitAction {
	var rateLimitActions []RateLimitAction
	result := initializers.DB.Find(&rateLimitActions, conditions, userId, action, timeFrom)
	if len(rateLimitActions) == 0 {
		return nil
	}
	if result.Error != nil {
		log.Print(fmt.Sprintf("Error finding rate limit actions of type: %s with user ID: %d: %s", action, userId, result.Error.Error()))
		return nil
	}
	return rateLimitActions
}

func GetRateLimitRecordsViaAlternativeIdentifier(conditions string, alternativeIdentifier string, action string, timeFrom time.Time) []RateLimitAction {
	var rateLimitActions []RateLimitAction
	result := initializers.DB.Find(&rateLimitActions, conditions, alternativeIdentifier, action, timeFrom)
	if len(rateLimitActions) == 0 {
		return nil
	}
	if result.Error != nil {
		log.Print(fmt.Sprintf("Error finding rate limit actions of type: %s with alternative identifier: %s: %s", action, alternativeIdentifier, result.Error.Error()))
		return nil
	}
	return rateLimitActions
}

func ClearOldRateLimitRecords(rateLimitActionType constants.RateLimitActionType, userId uint, alternativeIdentifier string) {
	timeFrom := utils.CurrentTimeMinusMinutesAsTime(rateLimitActionType.DurationInMinutes)
	var rateLimitActions []RateLimitAction
	if userId != 0 {
		rateLimitActions = GetRateLimitRecordsViaUserId("user_id = ? AND rate_limit_action = ? AND created_at < ?", userId, rateLimitActionType.Action, timeFrom)
	} else {
		rateLimitActions = GetRateLimitRecordsViaAlternativeIdentifier("alternative_identifier = ? AND rate_limit_action = ? AND created_at < ?", alternativeIdentifier, rateLimitActionType.Action, timeFrom)
	}

	if rateLimitActions == nil {
		return
	}
	result := initializers.DB.Delete(&rateLimitActions)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Error deleting old rate limit actions of type: %s with identifier: %d / %s: %s", rateLimitActionType.Action, userId, alternativeIdentifier, result.Error.Error()))
	}
}

func DeleteAllRateLimitRecordsForUser(tx *gorm.DB, userId uint) error {
	var rateLimitRecordsToDelete []RateLimitAction
	result := tx.Model(&RateLimitAction{}).Unscoped().Delete(&rateLimitRecordsToDelete, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Rate limit action deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
