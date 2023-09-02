package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
)

type workoutLink struct {
	WorkoutId       uint `json:"workoutId" binding:"required"`
	PositionInSplit int  `json:"positionInSplit" binding:"required"`
	Id              uint `json:"id"`
}

type templateSplit struct {
	WorkoutLinks []workoutLink `json:"workoutLinks" binding:"required"`
	Id           uint          `json:"id"`
	Description  string        `json:"description" binding:"required"`
	Duration     int           `json:"duration" binding:"required"`
}

type allTemplateSplits []templateSplit

func PutTemplateSplits(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var body allTemplateSplits
	err := c.ShouldBindJSON(&body)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var splitsToUpdateOrCreate []models.TemplateSplit
	for order, split := range body {
		splitsToUpdateOrCreate = append(splitsToUpdateOrCreate, models.TemplateSplitToUpdateOrCreate(userId.(uint), split.Description, split.Duration, split.Id, order))
	}

	tx := initializers.DB.Begin()

	savedSplitsIds, err := models.HandleTemplateSplitSave(tx, splitsToUpdateOrCreate, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	var splitWorkoutLinksToUpdateOrCreate []models.TemplateSplitWorkoutLink
	for i, split := range body {
		splitId := splitsToUpdateOrCreate[i].ID
		var linksForSplit []models.TemplateSplitWorkoutLink
		for _, link := range split.WorkoutLinks {
			linksForSplit = append(linksForSplit, models.TemplateSplitWorkoutLinkToUpdateOrCreate(userId.(uint), splitId, link.WorkoutId, link.PositionInSplit, link.Id))
		}
		splitWorkoutLinksToUpdateOrCreate = append(splitWorkoutLinksToUpdateOrCreate, linksForSplit...)
	}

	err = models.ValidateWorkoutIdsForSplit(splitWorkoutLinksToUpdateOrCreate, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	savedSplitWorkoutLinkIds, err := models.HandleTemplateSplitWorkoutLinkSave(tx, splitWorkoutLinksToUpdateOrCreate, userId.(uint))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	err = models.HandleTemplateSplitDelete(tx, userId.(uint), savedSplitsIds)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	err = models.HandleTemplateSplitWorkoutLinkDelete(tx, userId.(uint), savedSplitWorkoutLinkIds)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Template splits updated"})
}

func GetTemplateSplits(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var templateSplits []models.TemplateSplit
	initializers.DB.Find(&templateSplits, "user_id = ?", userId)

	var templateSplitWorkoutLinks []models.TemplateSplitWorkoutLink
	initializers.DB.Find(&templateSplitWorkoutLinks, "user_id = ?", userId)

	orderedSplits := make([]models.TemplateSplit, len(templateSplits))
	for _, split := range templateSplits {
		orderedSplits[split.Order] = split
	}

	var allTemplateSplits allTemplateSplits
	for _, split := range orderedSplits {
		var templateSplit = templateSplit{
			Id:          split.ID,
			Description: split.Description,
			Duration:    split.Duration,
		}

		var relevantWorkoutLinks []workoutLink
		for _, link := range templateSplitWorkoutLinks {
			if link.SplitId == split.ID {
				relevantWorkoutLinks = append(relevantWorkoutLinks, workoutLink{
					WorkoutId:       link.WorkoutId,
					PositionInSplit: link.PositionInSplit,
					Id:              link.ID,
				})
			}
		}

		sortedWorkoutLinks := relevantWorkoutLinks
		isSorted := false
		for !isSorted {
			isSorted = true
			i := 0
			for i < len(sortedWorkoutLinks)-1 {
				if sortedWorkoutLinks[i].PositionInSplit > sortedWorkoutLinks[i+1].PositionInSplit {
					sortedWorkoutLinks[i], sortedWorkoutLinks[i+1] = sortedWorkoutLinks[i+1], sortedWorkoutLinks[i]
					isSorted = false
				}
				i++
			}
		}

		templateSplit.WorkoutLinks = sortedWorkoutLinks
		allTemplateSplits = append(allTemplateSplits, templateSplit)
	}

	c.JSON(http.StatusOK, allTemplateSplits)
}
