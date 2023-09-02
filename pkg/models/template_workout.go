package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/utils"
)

type TemplateWorkout struct {
	gorm.Model
	UserId      uint   `gorm:"not null;uniqueIndex:idx_user_id_description;uniqueIndex:idx_user_id_order"`
	Order       int    `gorm:"not null;uniqueIndex:idx_user_id_order"`
	Description string `gorm:"not null;uniqueIndex:idx_user_id_description"`
}

func TemplateWorkoutToUpdateOrCreate(userId uint, orderInBundle int, description string, id uint) TemplateWorkout {
	return TemplateWorkout{
		UserId:      userId,
		Order:       orderInBundle,
		Description: description,
		Model:       gorm.Model{ID: id},
	}
}

func HandleTemplateWorkoutSave(tx *gorm.DB, workoutsToUpdateOrCreate []TemplateWorkout, userId uint) ([]uint, error) {
	result := tx.Model(&TemplateWorkout{}).Save(&workoutsToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template workout save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, result.Error
	}

	for _, workout := range workoutsToUpdateOrCreate {
		if workout.UserId != userId {
			return nil, errors.New("items within JSON do not belong to this user")
		}
	}

	savedWorkoutsIds := make([]uint, len(workoutsToUpdateOrCreate))
	for _, workout := range workoutsToUpdateOrCreate {
		savedWorkoutsIds[workout.Order] = workout.ID
	}
	return utils.RemoveZerosFromSliceOfUint(savedWorkoutsIds), nil
}

func HandleTemplateWorkoutDelete(tx *gorm.DB, savedWorkoutsIds []uint, userId uint) error {
	var workoutsToDelete []TemplateWorkout
	result := tx.Model(&TemplateWorkout{}).Unscoped().Delete(&workoutsToDelete, "user_id = ? AND id not in ?", userId, savedWorkoutsIds)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template workout deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
