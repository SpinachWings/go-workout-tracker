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

func CalendarSetToUpdateOrCreate(exerciseName string, userId uint, workoutId uint, orderInExercise int, weight int, kgOrLbs string, reps int, isometricHoldTime int, id uint, workoutDate string, exerciseIsIsometric bool, workoutIsCompleted bool) (CalendarSet, error) {
	repsToUse := reps
	if exerciseIsIsometric {
		repsToUse = 1
	}

	isometricHoldTimeToUse := isometricHoldTime
	if !exerciseIsIsometric {
		isometricHoldTimeToUse = 0
	}

	if workoutIsCompleted && exerciseIsIsometric && isometricHoldTime == 0 {
		return BlankCalendarSet(), errors.New("isometric hold time must be filled in for all sets of isometric exercises if workout is complete")
	}

	if workoutIsCompleted && !exerciseIsIsometric && reps == 0 {
		return BlankCalendarSet(), errors.New("reps must be filled in for all sets of relevant exercises if workout is complete")
	}

	if workoutIsCompleted && !exerciseIsIsometric && weight > 0 && kgOrLbs == "" {
		return BlankCalendarSet(), errors.New("kg or lbs must be defined for all sets of relevant exercises if workout is complete")
	}

	if kgOrLbs != "kg" && kgOrLbs != "lbs" && kgOrLbs != "" {
		return BlankCalendarSet(), errors.New("incorrect value for kg or lbs")
	}

	return CalendarSet{
		ExerciseName:      utils.StandardiseCase(exerciseName),
		UserId:            userId,
		WorkoutId:         workoutId,
		OrderInExercise:   orderInExercise,
		Weight:            weight,
		KgOrLbs:           kgOrLbs,
		Reps:              repsToUse,
		IsometricHoldTime: isometricHoldTimeToUse,
		WorkoutDate:       workoutDate,
		Model:             gorm.Model{ID: id},
	}, nil
}

func BlankCalendarSet() CalendarSet {
	return CalendarSet{
		ExerciseName:      "",
		UserId:            0,
		WorkoutId:         0,
		OrderInExercise:   0,
		Weight:            0,
		KgOrLbs:           "",
		Reps:              0,
		IsometricHoldTime: 0,
		WorkoutDate:       "",
		Model:             gorm.Model{ID: 0},
	}
}

func HandleCalendarSetSave(tx *gorm.DB, setsToUpdateOrCreate []CalendarSet, userId uint) ([]uint, error) {
	result := tx.Model(&CalendarSet{}).Save(&setsToUpdateOrCreate)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrEmptySlice) {
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

func HandleCalendarSetDelete(tx *gorm.DB, savedWorkoutsIds []uint, savedSetIds []uint, userId uint, condition string) error {
	var setsToDelete []CalendarSet
	result := tx.Model(&CalendarSet{}).Unscoped().Delete(&setsToDelete, condition, userId, savedWorkoutsIds, append(savedSetIds, 0))
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar set deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
