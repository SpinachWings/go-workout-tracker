package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/utils"
)

type TemplateExercise struct {
	gorm.Model
	ExerciseName     string `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	UserId           uint   `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	WorkoutId        uint   `gorm:"not null;uniqueIndex:idx_user_id_workout_id_exercise_name"`
	OrderInWorkout   int    `gorm:"not null"`
	IsIsometric      bool   `gorm:"not null;default:false"`
	IsPartOfSuperset bool   `gorm:"not null;default:false"`
	OrderInSuperset  int
}

func TemplateExerciseToUpdateOrCreate(exerciseName string, userId uint, workoutId uint, orderInWorkout int, isIsometric bool, orderInSuperset int, id uint) TemplateExercise {
	var isPartOfSuperset bool
	if orderInSuperset >= 0 {
		isPartOfSuperset = true
	}
	return TemplateExercise{
		ExerciseName:     utils.StandardiseCase(exerciseName),
		UserId:           userId,
		WorkoutId:        workoutId,
		OrderInWorkout:   orderInWorkout,
		IsIsometric:      isIsometric,
		IsPartOfSuperset: isPartOfSuperset,
		OrderInSuperset:  orderInSuperset,
		Model:            gorm.Model{ID: id},
	}
}

func HandleTemplateExerciseSave(tx *gorm.DB, exercisesToUpdateOrCreate []TemplateExercise, userId uint) ([]uint, error) {
	result := tx.Model(&TemplateExercise{}).Save(&exercisesToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template exercise save failed for user with ID: %d: %s", userId, result.Error.Error()))
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

func HandleTemplateExerciseDelete(tx *gorm.DB, savedExerciseIds []uint, userId uint) error {
	var exercisesToDelete []TemplateExercise
	result := tx.Model(&TemplateExercise{}).Unscoped().Delete(&exercisesToDelete, "user_id = ? AND id not in ?", userId, append(savedExerciseIds, 0))
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template exercise deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
