package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"workout-tracker-go-app/pkg/utils"
)

type TemplateSplit struct {
	gorm.Model
	UserId      uint   `gorm:"not null;uniqueIndex:idx_user_id_description;uniqueIndex:idx_user_id_order"`
	Description string `gorm:"not null;uniqueIndex:idx_user_id_description"`
	Duration    int    `gorm:"not null"`
	Order       int    `gorm:"not null;uniqueIndex:idx_user_id_order"`
}

func TemplateSplitToUpdateOrCreate(userId uint, description string, duration int, id uint, order int) TemplateSplit {
	return TemplateSplit{
		UserId:      userId,
		Description: description,
		Duration:    duration,
		Order:       order,
		Model:       gorm.Model{ID: id},
	}
}

func HandleTemplateSplitSave(tx *gorm.DB, splitsToUpdateOrCreate []TemplateSplit, userId uint) ([]uint, error) {
	result := tx.Model(TemplateSplit{}).Save(&splitsToUpdateOrCreate)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template split save failed for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, result.Error
	}

	savedSplitIds := make([]uint, len(splitsToUpdateOrCreate))
	for _, split := range splitsToUpdateOrCreate {
		savedSplitIds = append(savedSplitIds, split.ID)
	}
	return utils.RemoveZerosFromSliceOfUint(savedSplitIds), nil
}

func HandleTemplateSplitDelete(tx *gorm.DB, userId uint, savedSplitsIds []uint) error {
	var splitsToDelete TemplateSplit
	result := tx.Model(TemplateSplit{}).Unscoped().Delete(&splitsToDelete, "user_id = ? AND id not in ?", userId, savedSplitsIds)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Template split deletion failed for user with ID: %d: %s", userId, result.Error.Error()))
		return result.Error
	}
	return nil
}
