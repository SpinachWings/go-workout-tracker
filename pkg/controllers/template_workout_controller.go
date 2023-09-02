package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
)

type templateWorkout struct {
	Exercises   []services.TemplateExerciseGroup `json:"exercises"`
	Id          uint                             `json:"id"`
	Description string                           `json:"description" binding:"required"`
}

type allTemplateWorkouts []templateWorkout

func PutTemplateWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	//userId = userId.(uint) + 1 // just for testing crap
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var body allTemplateWorkouts
	err := c.ShouldBindJSON(&body)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var workoutsToUpdateOrCreate []models.TemplateWorkout
	for orderInBundle, workout := range body {
		workoutsToUpdateOrCreate = append(workoutsToUpdateOrCreate, models.TemplateWorkoutToUpdateOrCreate(userId.(uint), orderInBundle, workout.Description, workout.Id))
	}

	tx := initializers.DB.Begin()

	savedWorkoutsIds, err := models.HandleTemplateWorkoutSave(tx, workoutsToUpdateOrCreate, userId.(uint))
	if err != nil && err.Error() == "items within JSON do not belong to this user" {
		tx.Rollback()
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	var exercisesToUpdateOrCreate []models.TemplateExercise
	for orderInBundle, workout := range body {
		workoutId := savedWorkoutsIds[orderInBundle]
		for orderInWorkout, exerciseGroup := range workout.Exercises {
			exerciseGroupConverter := services.TemplateExerciseGroupConverter{
				ExerciseGroup:  exerciseGroup,
				UserId:         userId.(uint),
				WorkoutId:      workoutId,
				OrderInWorkout: orderInWorkout,
			}
			exerciseGroupConverter.TemplateExerciseGroupToRelevantModel(&exercisesToUpdateOrCreate)
		}
	}

	savedExerciseIds, err := models.HandleTemplateExerciseSave(tx, exercisesToUpdateOrCreate, userId.(uint))
	if err != nil && err.Error() == "items within JSON do not belong to this user" {
		tx.Rollback()
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	err = models.HandleTemplateWorkoutDelete(tx, savedWorkoutsIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	err = models.HandleTemplateExerciseDelete(tx, savedExerciseIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Template workouts updated"})
}

func GetTemplateWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var templateWorkouts []models.TemplateWorkout
	initializers.DB.Find(&templateWorkouts, "user_id = ?", userId)

	var templateExercises []models.TemplateExercise
	initializers.DB.Find(&templateExercises, "user_id = ?", userId)

	orderedWorkouts := make([]models.TemplateWorkout, len(templateWorkouts))
	for _, workout := range templateWorkouts {
		orderedWorkouts[workout.Order] = workout
	}

	var allTemplateWorkouts allTemplateWorkouts

	for _, workout := range orderedWorkouts {
		relevantExercises := services.GetTemplateExercisesInThisWorkout(templateExercises, workout.ID)
		numberOfExerciseGroups := services.GetNumberOfTemplateExerciseGroupsInThisWorkout(relevantExercises)
		orderedExerciseGroups := services.GetTemplateExerciseGroupsInThisWorkout(numberOfExerciseGroups, relevantExercises)

		templateWorkout := templateWorkout{
			Id:          workout.ID,
			Description: workout.Description,
			Exercises:   orderedExerciseGroups,
		}
		allTemplateWorkouts = append(allTemplateWorkouts, templateWorkout)
	}

	c.JSON(http.StatusOK, allTemplateWorkouts)
}
