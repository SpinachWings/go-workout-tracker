package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
	"workout-tracker-go-app/pkg/utils"
)

type templateWorkout struct {
	Exercises   []services.TemplateExerciseGroup `json:"exercises"`
	Id          uint                             `json:"id"`
	Description string                           `json:"description" binding:"required"`
}

type allTemplateWorkouts []templateWorkout

// templates - save, update & delete EVERYTHING for the user at once

func PutTemplateWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var body allTemplateWorkouts
	err := c.ShouldBindJSON(&body)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var workoutsToUpdateOrCreate []models.TemplateWorkout
	for orderInBundle, workout := range body {
		workoutsToUpdateOrCreate = append(workoutsToUpdateOrCreate, models.TemplateWorkoutToUpdateOrCreate(userId.(uint), orderInBundle, workout.Description, workout.Id))
	}
	maxWorkouts := constants.GetRestrictions().TemplateWorkoutsPerUser.GetRestrictionAmount(false)
	if len(workoutsToUpdateOrCreate) > maxWorkouts {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("number of workouts cannot exceed: %d", maxWorkouts)})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	var exercisesToUpdateOrCreate []models.TemplateExercise
	for orderInBundle, workout := range body {
		exerciseGroupsPerWorkoutMaxLength := constants.GetRestrictions().TemplateMaxExerciseGroupsPerWorkout.GetRestrictionAmount(false)
		if len(workout.Exercises) > exerciseGroupsPerWorkoutMaxLength {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("number of exercises grousp per workout cannot exceed: %d", exerciseGroupsPerWorkoutMaxLength)})
			return
		}

		workoutId := savedWorkoutsIds[orderInBundle]
		for orderInWorkout, exerciseGroup := range workout.Exercises {
			exercisesPerGroupMaxLength := constants.GetRestrictions().TemplateMaxExercisesPerGroup.GetRestrictionAmount(false)
			if len(exerciseGroup) > exercisesPerGroupMaxLength {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("number of exercises per group (superset) cannot exceed: %d", exercisesPerGroupMaxLength)})
				return
			}

			exerciseGroupConverter := services.TemplateExerciseGroupConverter{
				ExerciseGroup:          exerciseGroup,
				UserId:                 userId.(uint),
				WorkoutId:              workoutId,
				OrderInWorkout:         orderInWorkout,
				OrderOfWorkoutInBundle: orderInBundle,
			}
			exerciseGroupConverter.TemplateExerciseGroupToRelevantModel(&exercisesToUpdateOrCreate)
		}
	}

	var allExerciseNamesWithOrderOfWorkoutInBundle []string
	for _, exercise := range exercisesToUpdateOrCreate {
		allExerciseNamesWithOrderOfWorkoutInBundle = append(allExerciseNamesWithOrderOfWorkoutInBundle, fmt.Sprintf("%s_%d", exercise.ExerciseName, exercise.OrderOfWorkoutInBundle))
	}
	if utils.SliceOfStringContainsDuplicates(allExerciseNamesWithOrderOfWorkoutInBundle) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "the same exercise cannot appear twice in the same workout"})
		return
	}

	savedExerciseIds, err := models.HandleTemplateExerciseSave(tx, exercisesToUpdateOrCreate, userId.(uint))
	if err != nil && err.Error() == "items within JSON do not belong to this user" {
		tx.Rollback()
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	err = models.HandleTemplateWorkoutDelete(tx, savedWorkoutsIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	err = models.HandleTemplateExerciseDelete(tx, savedExerciseIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "template workouts updated"})
}

func GetTemplateWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	templateWorkouts, templateExercises, err := services.FindTemplateWorkouts(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

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
