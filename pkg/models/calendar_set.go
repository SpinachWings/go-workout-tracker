package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/utils"
)

type CalendarSet struct {
	gorm.Model
	UserId            uint   `gorm:"not null"`
	WorkoutId         uint   `gorm:"not null"`
	ExerciseName      string `gorm:"not null"`
	Weight            int
	KgOrLbs           string
	Reps              int
	IsometricHoldTime int
	OrderInExercise   int `gorm:"not null"`
	WorkoutDate       string
}

// if isIsometric - reps becomes 1
// isometric hold time cannot be null if is isometric
// reps OR isometricHoldTime, weight & kgorlbs cannot be null if workout id is completed
// ensure kgOrLbs is kg or lbs

func CalendarSetToUpdateOrCreate(exerciseName string, userId uint, workoutId uint, orderInExercise int, weight int, kgOrLbs string, reps int, isometricHoldTime int, id uint, workoutDate string) CalendarSet {
	return CalendarSet{
		ExerciseName:      exerciseName,
		UserId:            userId,
		WorkoutId:         workoutId,
		OrderInExercise:   orderInExercise,
		Weight:            weight,
		KgOrLbs:           kgOrLbs,
		Reps:              reps,
		IsometricHoldTime: isometricHoldTime,
		WorkoutDate:       workoutDate,
		Model:             gorm.Model{ID: id},
	}
}

func HandleCalendarSetSave(tx *gorm.DB, setsToUpdateOrCreate []CalendarSet, userId uint) ([]uint, error) {
	result := tx.Model(&CalendarSet{}).Save(&setsToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar set save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, result.Error
	}

	for _, set := range setsToUpdateOrCreate {
		if set.UserId != userId {
			return nil, errors.New("items within JSON do not belong to this user")
		}
	}

	savedSetIds := make([]uint, len(setsToUpdateOrCreate))
	for _, set := range setsToUpdateOrCreate {
		savedSetIds = append(savedSetIds, set.ID)
	}

	return utils.RemoveZerosFromSliceOfUint(savedSetIds), nil
}

func HandleCalendarSetDelete(tx *gorm.DB, savedWorkoutsIds []uint, savedSetIds []uint, userId uint) error {
	var setsToDelete []CalendarSet
	result := tx.Model(&CalendarSet{}).Unscoped().Delete(&setsToDelete, "user_id = ? AND workout_id in ? AND id not in ?", userId, savedWorkoutsIds, savedSetIds)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar set deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
