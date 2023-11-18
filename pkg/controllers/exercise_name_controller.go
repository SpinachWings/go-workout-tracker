package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/utils"
)

func GetExerciseNames(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var templateExerciseNames []string
	result := initializers.DB.Model(&models.TemplateExercise{}).Where("user_id = ?", userId.(uint)).Distinct().Pluck("exercise_name", &templateExerciseNames)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding template exercises for exercise name dropdown list for user with ID: %d: %s", userId, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	var calendarExerciseNames []string
	result = initializers.DB.Model(&models.CalendarExercise{}).Where("user_id = ?", userId.(uint)).Distinct().Pluck("exercise_name", &calendarExerciseNames)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding calendar exercises for exercise name dropdown list for user with ID: %d: %s", userId, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	c.JSON(http.StatusOK, utils.CombineTwoSlicesOfStringNoDuplicates(templateExerciseNames, calendarExerciseNames))
}
