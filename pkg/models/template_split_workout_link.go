package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/utils"
)

type TemplateSplitWorkoutLink struct {
	gorm.Model
	UserId          uint `gorm:"not null;uniqueIndex:idx_user_id_split_id_workout_id;uniqueIndex:idx_user_id_split_id_position"`
	SplitId         uint `gorm:"not null;uniqueIndex:idx_user_id_split_id_workout_id;uniqueIndex:idx_user_id_split_id_position"`
	WorkoutId       uint `gorm:"not null;uniqueIndex:idx_user_id_split_id_workout_id"`
	PositionInSplit int  `gorm:"not null;uniqueIndex:idx_user_id_split_id_position"`
}

func TemplateSplitWorkoutLinkToUpdateOrCreate(userId uint, splitId uint, workoutId uint, positionInSplit int, id uint) TemplateSplitWorkoutLink {
	return TemplateSplitWorkoutLink{
		UserId:          userId,
		SplitId:         splitId,
		WorkoutId:       workoutId,
		PositionInSplit: positionInSplit,
		Model:           gorm.Model{ID: id},
	}
}

func ValidateWorkoutIdsForSplit(splitWorkoutLinksToUpdateOrCreate []TemplateSplitWorkoutLink, userId uint) error {
	var workoutIds []uint
	for _, link := range splitWorkoutLinksToUpdateOrCreate {
		workoutIds = append(workoutIds, link.WorkoutId)
	}

	var templateWorkouts []TemplateWorkout
	result := initializers.DB.Find(&templateWorkouts, "id in ?", workoutIds)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding template workouts for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}

	if len(templateWorkouts) != len(workoutIds) {
		return errors.New("all workout ids must exist and belong to the user")
	}

	for _, workout := range templateWorkouts {
		if workout.UserId != userId {
			return errors.New("all workout ids must exist and belong to the user")
		}
	}

	return nil
}

func HandleTemplateSplitWorkoutLinkSave(tx *gorm.DB, splitWorkoutLinksToUpdateOrCreate []TemplateSplitWorkoutLink, userId uint) ([]uint, error) {
	result := tx.Model(TemplateSplitWorkoutLink{}).Save(&splitWorkoutLinksToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template split workout link save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, result.Error
	}

	for _, link := range splitWorkoutLinksToUpdateOrCreate {
		if link.UserId != userId {
			return nil, errors.New("items within JSON do not belong to this user")
		}
	}

	savedSplitWorkoutLinkIds := make([]uint, len(splitWorkoutLinksToUpdateOrCreate))
	for _, link := range splitWorkoutLinksToUpdateOrCreate {
		savedSplitWorkoutLinkIds = append(savedSplitWorkoutLinkIds, link.ID)
	}
	return utils.RemoveZerosFromSliceOfUint(savedSplitWorkoutLinkIds), nil
}

func HandleTemplateSplitWorkoutLinkDelete(tx *gorm.DB, userId uint, savedSplitWorkoutLinkIds []uint) error {
	var splitWorkoutLinksToDelete TemplateSplitWorkoutLink
	result := tx.Model(TemplateSplitWorkoutLink{}).Unscoped().Delete(&splitWorkoutLinksToDelete, "user_id = ? AND id not in ?", userId, append(savedSplitWorkoutLinkIds, 0))
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template split deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
