package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"regexp"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/utils"
)

type CalendarWorkout struct {
	gorm.Model
	UserId      uint   `gorm:"not null;uniqueIndex:idx_user_id_date"`
	Date        string `gorm:"not null;uniqueIndex:idx_user_id_date"`
	Description string
	IsCompleted bool `gorm:"not null;default:false"`
}

func CalendarWorkoutToUpdateOrCreate(userId uint, date string, description string, isCompleted bool, id uint) (CalendarWorkout, error) {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	if !re.MatchString(date) {
		return BlankCalendarWorkout(), errors.New("date format for workout must be yyyy-mm-dd")
	}

	if utils.DateAsStringIsInFuture(date) && isCompleted {
		return BlankCalendarWorkout(), errors.New("workout cannot be completed if date is in future")
	}

	maxYears := constants.GetRestrictions().CalendarWorkoutMaxDateYears.GetRestrictionAmount(false)
	if utils.DateAsStringIsMoreThanNumYearsInFuture(date, maxYears) {
		return BlankCalendarWorkout(), errors.New(fmt.Sprintf("date cannot be more than: %d years in the future", maxYears))
	}

	minYears := constants.GetRestrictions().CalendarWorkoutMinDateYears.GetRestrictionAmount(false)
	if utils.DateAsStringIsLessThanNumYearsInPast(date, minYears) {
		return BlankCalendarWorkout(), errors.New(fmt.Sprintf("date cannot be less than: %d years in the past", minYears))
	}

	return CalendarWorkout{
		UserId:      userId,
		Date:        date,
		Description: description,
		IsCompleted: isCompleted,
		Model:       gorm.Model{ID: id},
	}, nil
}

func BlankCalendarWorkout() CalendarWorkout {
	return CalendarWorkout{
		UserId:      0,
		Date:        "",
		Description: "",
		IsCompleted: false,
	}
}

func HandleCalendarWorkoutSave(tx *gorm.DB, workoutsToUpdateOrCreate []CalendarWorkout, userId uint) (map[string]uint, []uint, error) {
	result := tx.Model(&CalendarWorkout{}).Save(&workoutsToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar workout save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, result.Error
	}

	for _, workout := range workoutsToUpdateOrCreate {
		if workout.UserId != userId {
			return nil, nil, errors.New("items within JSON do not belong to this user")
		}
	}

	savedWorkoutsIdsMap := make(map[string]uint, len(workoutsToUpdateOrCreate))
	for _, workout := range workoutsToUpdateOrCreate {
		savedWorkoutsIdsMap[workout.Date] = workout.ID
	}
	savedWorkoutsIdsSlice := make([]uint, len(workoutsToUpdateOrCreate))
	for _, workout := range workoutsToUpdateOrCreate {
		savedWorkoutsIdsSlice = append(savedWorkoutsIdsSlice, workout.ID)
	}
	return savedWorkoutsIdsMap, utils.RemoveZerosFromSliceOfUint(savedWorkoutsIdsSlice), nil
}

func HandleCalendarWorkoutDelete(tx *gorm.DB, workoutIdsToDelete []uint, userId uint, condition string) error {
	var workoutsToDelete []CalendarWorkout
	result := tx.Model(&CalendarWorkout{}).Unscoped().Delete(&workoutsToDelete, condition, userId, workoutIdsToDelete)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar workout deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
