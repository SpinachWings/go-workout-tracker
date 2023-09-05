package services

import (
	"fmt"
	"log"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
)

func FindTemplateWorkouts(userId uint) ([]models.TemplateWorkout, []models.TemplateExercise, error) {
	var templateWorkouts []models.TemplateWorkout
	result := initializers.DB.Find(&templateWorkouts, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find template workouts for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, result.Error
	}

	var templateExercises []models.TemplateExercise
	result = initializers.DB.Find(&templateExercises, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find template exercises for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, result.Error
	}

	return templateWorkouts, templateExercises, nil
}

func FindTemplateSplits(userId uint) ([]models.TemplateSplit, []models.TemplateSplitWorkoutLink, error) {
	var templateSplits []models.TemplateSplit
	result := initializers.DB.Find(&templateSplits, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find template splits for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, result.Error
	}

	var templateSplitWorkoutLinks []models.TemplateSplitWorkoutLink
	result = initializers.DB.Find(&templateSplitWorkoutLinks, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find template split workout links for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, result.Error
	}

	return templateSplits, templateSplitWorkoutLinks, nil
}

func FindCalendarWorkouts(userId uint) ([]models.CalendarWorkout, []models.CalendarExercise, []models.CalendarSet, error) {
	var calendarWorkouts []models.CalendarWorkout
	result := initializers.DB.Find(&calendarWorkouts, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find calendar workouts for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, nil, result.Error
	}

	var calendarExercises []models.CalendarExercise
	result = initializers.DB.Find(&calendarExercises, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find calendar exercises for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, nil, result.Error
	}

	var calendarSets []models.CalendarSet
	result = initializers.DB.Find(&calendarSets, "user_id = ?", userId)
	if result.Error != nil {
		log.Print(fmt.Sprintf("Failed to find calendar sets for user with ID: %d: %s", userId, result.Error.Error()))
		return nil, nil, nil, result.Error
	}

	return calendarWorkouts, calendarExercises, calendarSets, nil
}
