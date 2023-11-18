package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/utils"
)

type CalendarExercise struct {
	gorm.Model
	ExerciseName     string `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	UserId           uint   `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	WorkoutId        uint   `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	OrderInWorkout   int    `gorm:"not null"`
	IsIsometric      bool   `gorm:"not null;default:false"`
	IsPartOfSuperset bool   `gorm:"not null;default:false"`
	OrderInSuperset  int
	WorkoutDate      string
}

func CalendarExerciseToUpdateOrCreate(exerciseName string, userId uint, workoutId uint, orderInWorkout int, isIsometric bool, orderInSuperset int, id uint, workoutDate string) CalendarExercise {
	var isPartOfSuperset bool
	if orderInSuperset >= 0 {
		isPartOfSuperset = true
	}
	return CalendarExercise{
		ExerciseName:     utils.StandardiseCase(exerciseName),
		UserId:           userId,
		WorkoutId:        workoutId,
		OrderInWorkout:   orderInWorkout,
		IsIsometric:      isIsometric,
		IsPartOfSuperset: isPartOfSuperset,
		OrderInSuperset:  orderInSuperset,
		WorkoutDate:      workoutDate,
		Model:            gorm.Model{ID: id},
	}
}

func HandleCalendarExerciseSave(tx *gorm.DB, exercisesToUpdateOrCreate []CalendarExercise, userId uint) ([]uint, error) {
	result := tx.Model(&CalendarExercise{}).Save(&exercisesToUpdateOrCreate)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrEmptySlice) {
		log.Print(fmt.Sprintf("Calendar exercise save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, result.Error
	}

	for _, exercise := range exercisesToUpdateOrCreate {
		if exercise.UserId != userId {
			return nil, errors.New("items within JSON do not belong to this user")
		}
	}

	savedExerciseIds := make([]uint, len(exercisesToUpdateOrCreate))
	for _, exercise := range exercisesToUpdateOrCreate {
		savedExerciseIds = append(savedExerciseIds, exercise.ID)
	}
	return utils.RemoveZerosFromSliceOfUint(savedExerciseIds), nil
}

func HandleCalendarExerciseDelete(tx *gorm.DB, savedWorkoutsIds []uint, savedExerciseIds []uint, userId uint) error {
	var exercisesToDelete []CalendarExercise
	result := tx.Model(&CalendarExercise{}).Unscoped().Delete(&exercisesToDelete, "user_id = ? AND workout_id in ? AND id not in ?", userId, savedWorkoutsIds, append(savedExerciseIds, 0))
	if result.Error != nil {
		log.Print(fmt.Sprintf("Calendar exercise deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
