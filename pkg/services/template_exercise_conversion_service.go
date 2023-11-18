package services

import (
	"workout-tracker-go-app/pkg/models"
)

type TemplateExercise struct {
	ExerciseName string `json:"exerciseName" binding:"required"`
	IsIsometric  bool   `json:"isIsometric" binding:"required"`
	Id           uint   `json:"id"`
}

type TemplateExerciseGroup []TemplateExercise

type TemplateExerciseGroupConverter struct {
	ExerciseGroup          TemplateExerciseGroup
	UserId                 uint
	WorkoutId              uint
	OrderInWorkout         int
	OrderOfWorkoutInBundle int
}

func (egc TemplateExerciseGroupConverter) TemplateExerciseGroupToRelevantModel(exercisesToUpdateOrCreate *[]models.TemplateExercise) {
	if len(egc.ExerciseGroup) >= 2 {
		egc.ConvertTemplateSuperset(exercisesToUpdateOrCreate)
	} else {
		egc.ConvertTemplateSingleExercise(exercisesToUpdateOrCreate)
	}
}

func (egc TemplateExerciseGroupConverter) ConvertTemplateSuperset(exercisesToUpdateOrCreate *[]models.TemplateExercise) {
	for orderInSuperset, exercise := range egc.ExerciseGroup {
		*exercisesToUpdateOrCreate = append(*exercisesToUpdateOrCreate, models.TemplateExerciseToUpdateOrCreate(exercise.ExerciseName, egc.UserId, egc.WorkoutId, egc.OrderInWorkout, exercise.IsIsometric, orderInSuperset, exercise.Id, egc.OrderOfWorkoutInBundle))
	}
}

func (egc TemplateExerciseGroupConverter) ConvertTemplateSingleExercise(exercisesToUpdateOrCreate *[]models.TemplateExercise) {
	exercise := egc.ExerciseGroup[0]
	*exercisesToUpdateOrCreate = append(*exercisesToUpdateOrCreate, models.TemplateExerciseToUpdateOrCreate(exercise.ExerciseName, egc.UserId, egc.WorkoutId, egc.OrderInWorkout, exercise.IsIsometric, -1, exercise.Id, egc.OrderOfWorkoutInBundle))
}

func GetTemplateExercisesInThisWorkout(allTemplateExercises []models.TemplateExercise, workoutId uint) []models.TemplateExercise {
	var relevantExercises []models.TemplateExercise
	for _, exercise := range allTemplateExercises {
		if exercise.WorkoutId == workoutId {
			relevantExercises = append(relevantExercises, exercise)
		}
	}
	return relevantExercises
}

func GetNumberOfTemplateExerciseGroupsInThisWorkout(relevantExercises []models.TemplateExercise) int {
	numberOfExerciseGroups := 0
	for _, exercise := range relevantExercises {
		if exercise.OrderInWorkout+1 > numberOfExerciseGroups {
			numberOfExerciseGroups = exercise.OrderInWorkout + 1
		}
	}
	return numberOfExerciseGroups
}

func GetTemplateExerciseGroupsInThisWorkout(numberOfExerciseGroups int, relevantExercises []models.TemplateExercise) []TemplateExerciseGroup {
	orderedExerciseGroups := make([]TemplateExerciseGroup, numberOfExerciseGroups)
	for orderInWorkout, exercise := range relevantExercises {
		if orderedExerciseGroups[exercise.OrderInWorkout] == nil {
			exercisesAtThisOrder := GetTemplateExercisesInThisExerciseGroup(relevantExercises, orderInWorkout)
			exerciseGroup := RelevantModelToTemplateExerciseGroup(exercisesAtThisOrder)
			orderedExerciseGroups[exercise.OrderInWorkout] = exerciseGroup
		}
	}
	return orderedExerciseGroups
}

func GetTemplateExercisesInThisExerciseGroup(relevantExercises []models.TemplateExercise, orderInWorkout int) []models.TemplateExercise {
	var exercisesAtThisOrder []models.TemplateExercise
	for _, exercise := range relevantExercises {
		if exercise.OrderInWorkout == orderInWorkout {
			exercisesAtThisOrder = append(exercisesAtThisOrder, exercise)
		}
	}
	return exercisesAtThisOrder
}

func RelevantModelToTemplateExerciseGroup(exercisesAtThisOrder []models.TemplateExercise) TemplateExerciseGroup {
	var exerciseGroup TemplateExerciseGroup
	if len(exercisesAtThisOrder) >= 2 {
		exerciseGroup = ConvertToTemplateSuperset(exercisesAtThisOrder)
	} else {
		exerciseGroup = ConvertToTemplateSingleExercise(exercisesAtThisOrder)
	}
	return exerciseGroup
}

func ConvertToTemplateSuperset(exercisesAtThisOrder []models.TemplateExercise) TemplateExerciseGroup {
	var exerciseGroup TemplateExerciseGroup
	orderedSupersetExercises := make([]models.TemplateExercise, len(exercisesAtThisOrder))
	for _, exercise := range exercisesAtThisOrder {
		orderedSupersetExercises[exercise.OrderInSuperset] = exercise
	}

	for _, exercise := range orderedSupersetExercises {
		exerciseGroup = append(exerciseGroup, TemplateExercise{
			ExerciseName: exercise.ExerciseName,
			IsIsometric:  exercise.IsIsometric,
			Id:           exercise.ID,
		})
	}
	return exerciseGroup
}

func ConvertToTemplateSingleExercise(exercisesAtThisOrder []models.TemplateExercise) TemplateExerciseGroup {
	return TemplateExerciseGroup{
		TemplateExercise{
			ExerciseName: exercisesAtThisOrder[0].ExerciseName,
			IsIsometric:  exercisesAtThisOrder[0].IsIsometric,
			Id:           exercisesAtThisOrder[0].ID,
		},
	}
}
