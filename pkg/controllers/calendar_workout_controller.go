package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
)

type calendarWorkout struct {
	Exercises   []services.CalendarExerciseGroup `json:"exercises"`
	Id          uint                             `json:"id"`
	Date        string                           `json:"date" binding:"required"`
	Description string                           `json:"description" binding:"required"`
	IsCompleted bool                             `json:"isCompleted" binding:"required"`
}

type allCalendarWorkouts []calendarWorkout

type allCalendarWorkoutsIncludingItemsToDelete struct {
	AllCalendarWorkouts allCalendarWorkouts `json:"allCalendarWorkouts"`
	WorkoutIdsToDelete  []uint              `json:"workoutIdsToDelete"`
}

func PutCalendarWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var body allCalendarWorkoutsIncludingItemsToDelete
	err := c.ShouldBindJSON(&body)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tx := initializers.DB.Begin()

	err = models.HandleCalendarWorkoutDelete(tx, body.WorkoutIdsToDelete, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	var workoutsToUpdateOrCreate []models.CalendarWorkout
	for _, workout := range body.AllCalendarWorkouts {
		calendarWorkoutToUpdateOrCreate, err := models.CalendarWorkoutToUpdateOrCreate(userId.(uint), workout.Date, workout.Description, workout.IsCompleted, workout.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		workoutsToUpdateOrCreate = append(workoutsToUpdateOrCreate, calendarWorkoutToUpdateOrCreate)
	}

	savedWorkoutsIdsMap, savedWorkoutIdsSlice, err := models.HandleCalendarWorkoutSave(tx, workoutsToUpdateOrCreate, userId.(uint))
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

	var exercisesToUpdateOrCreate []models.CalendarExercise
	var setsToUpdateOrCreate []models.CalendarSet
	for _, workout := range body.AllCalendarWorkouts {
		workoutId := savedWorkoutsIdsMap[workout.Date]
		for orderInWorkout, exerciseGroup := range workout.Exercises {
			exerciseGroupConverter := services.CalendarExerciseGroupConverter{
				ExerciseGroup:  exerciseGroup,
				UserId:         userId.(uint),
				WorkoutId:      workoutId,
				OrderInWorkout: orderInWorkout,
				WorkoutDate:    workout.Date,
			}
			exerciseGroupConverter.CalendarExerciseGroupToRelevantModel(&exercisesToUpdateOrCreate, &setsToUpdateOrCreate)
		}
	}

	savedExerciseIds, err := models.HandleCalendarExerciseSave(tx, exercisesToUpdateOrCreate, userId.(uint))
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

	savedSetIds, err := models.HandleCalendarSetSave(tx, setsToUpdateOrCreate, userId.(uint))
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

	err = models.HandleCalendarExerciseDelete(tx, savedWorkoutIdsSlice, savedExerciseIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	err = models.HandleCalendarSetDelete(tx, savedWorkoutIdsSlice, savedSetIds, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Calendar workouts updated"})
}

func GetCalendarWorkouts(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	calendarWorkouts, calendarExercises, calendarSets, err := services.FindCalendarWorkouts(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	workoutsByDate := make(map[string]models.CalendarWorkout, len(calendarWorkouts))
	for _, workout := range calendarWorkouts {
		workoutsByDate[workout.Date] = workout
	}

	var allCalendarWorkouts allCalendarWorkouts

	for date, workout := range workoutsByDate {
		relevantExercises := services.GetCalendarExercisesInThisWorkout(calendarExercises, workout.ID)
		numberOfExerciseGroups := services.GetNumberOfCalendarExerciseGroupsInThisWorkout(relevantExercises)
		orderedExerciseGroups := services.GetCalendarExerciseGroupsInThisWorkout(numberOfExerciseGroups, relevantExercises, calendarSets, workout.ID)

		calendarWorkout := calendarWorkout{
			Id:          workout.ID,
			Date:        date,
			Description: workout.Description,
			IsCompleted: workout.IsCompleted,
			Exercises:   orderedExerciseGroups,
		}
		allCalendarWorkouts = append(allCalendarWorkouts, calendarWorkout)
	}

	c.JSON(http.StatusOK, allCalendarWorkouts)
}
