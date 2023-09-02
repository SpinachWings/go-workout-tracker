package services

import (
	"workout-tracker-go-app/pkg/models"
)

type Set struct {
	Weight            int
	KgOrLbs           string
	Reps              int
	IsometricHoldTime int
	Id                uint
}

type Sets []Set

type CalendarExercise struct {
	ExerciseName string `json:"exerciseName" binding:"required"`
	IsIsometric  bool   `json:"isIsometric" binding:"required"`
	Id           uint   `json:"id"`
	Sets         Sets   `json:"sets"`
}

type CalendarExerciseGroup []CalendarExercise

type CalendarExerciseGroupConverter struct {
	ExerciseGroup  CalendarExerciseGroup
	UserId         uint
	WorkoutId      uint
	OrderInWorkout int
	WorkoutDate    string
}

func (egc CalendarExerciseGroupConverter) CalendarExerciseGroupToRelevantModel(exercisesToUpdateOrCreate *[]models.CalendarExercise, setsToUpdateOrCreate *[]models.CalendarSet) {
	if len(egc.ExerciseGroup) >= 2 {
		egc.ConvertCalendarSuperset(exercisesToUpdateOrCreate, setsToUpdateOrCreate)
	} else {
		egc.ConvertCalendarSingleExercise(exercisesToUpdateOrCreate, setsToUpdateOrCreate)
	}
}

func (egc CalendarExerciseGroupConverter) ConvertCalendarSuperset(exercisesToUpdateOrCreate *[]models.CalendarExercise, setsToUpdateOrCreate *[]models.CalendarSet) {
	for orderInSuperset, exercise := range egc.ExerciseGroup {
		*exercisesToUpdateOrCreate = append(*exercisesToUpdateOrCreate, models.CalendarExerciseToUpdateOrCreate(exercise.ExerciseName, egc.UserId, egc.WorkoutId, egc.OrderInWorkout, exercise.IsIsometric, orderInSuperset, exercise.Id, egc.WorkoutDate))
		egc.ConvertSets(exercise, setsToUpdateOrCreate)
	}
}

func (egc CalendarExerciseGroupConverter) ConvertCalendarSingleExercise(exercisesToUpdateOrCreate *[]models.CalendarExercise, setsToUpdateOrCreate *[]models.CalendarSet) {
	exercise := egc.ExerciseGroup[0]
	*exercisesToUpdateOrCreate = append(*exercisesToUpdateOrCreate, models.CalendarExerciseToUpdateOrCreate(exercise.ExerciseName, egc.UserId, egc.WorkoutId, egc.OrderInWorkout, exercise.IsIsometric, -1, exercise.Id, egc.WorkoutDate))
	egc.ConvertSets(exercise, setsToUpdateOrCreate)
}

func (egc CalendarExerciseGroupConverter) ConvertSets(exercise CalendarExercise, setsToUpdateOrCreate *[]models.CalendarSet) {
	var setsForThisExercise []models.CalendarSet
	for orderInExercise, set := range exercise.Sets {
		setsForThisExercise = append(setsForThisExercise, models.CalendarSetToUpdateOrCreate(exercise.ExerciseName, egc.UserId, egc.WorkoutId, orderInExercise, set.Weight, set.KgOrLbs, set.Reps, set.IsometricHoldTime, set.Id, egc.WorkoutDate))
	}
	*setsToUpdateOrCreate = append(*setsToUpdateOrCreate, setsForThisExercise...)
}

func GetCalendarExercisesInThisWorkout(allCalendarExercises []models.CalendarExercise, workoutId uint) []models.CalendarExercise {
	var relevantExercises []models.CalendarExercise
	for _, exercise := range allCalendarExercises {
		if exercise.WorkoutId == workoutId {
			relevantExercises = append(relevantExercises, exercise)
		}
	}
	return relevantExercises
}

func GetNumberOfCalendarExerciseGroupsInThisWorkout(relevantExercises []models.CalendarExercise) int {
	numberOfExerciseGroups := 0
	for _, exercise := range relevantExercises {
		if exercise.OrderInWorkout+1 > numberOfExerciseGroups {
			numberOfExerciseGroups = exercise.OrderInWorkout + 1
		}
	}
	return numberOfExerciseGroups
}

func GetCalendarExerciseGroupsInThisWorkout(numberOfExerciseGroups int, relevantExercises []models.CalendarExercise, calendarSets []models.CalendarSet, workoutId uint) []CalendarExerciseGroup {
	orderedExerciseGroups := make([]CalendarExerciseGroup, numberOfExerciseGroups)
	for orderInWorkout, exercise := range relevantExercises {
		if orderedExerciseGroups[exercise.OrderInWorkout] == nil {
			exercisesAtThisOrder := GetCalendarExercisesInThisExerciseGroup(relevantExercises, orderInWorkout)
			exerciseGroup := RelevantModelToCalendarExerciseGroup(exercisesAtThisOrder, calendarSets, workoutId)
			orderedExerciseGroups[exercise.OrderInWorkout] = exerciseGroup
		}
	}
	return orderedExerciseGroups
}

func GetCalendarExercisesInThisExerciseGroup(relevantExercises []models.CalendarExercise, orderInWorkout int) []models.CalendarExercise {
	var exercisesAtThisOrder []models.CalendarExercise
	for _, exercise := range relevantExercises {
		if exercise.OrderInWorkout == orderInWorkout {
			exercisesAtThisOrder = append(exercisesAtThisOrder, exercise)
		}
	}
	return exercisesAtThisOrder
}

func RelevantModelToCalendarExerciseGroup(exercisesAtThisOrder []models.CalendarExercise, calendarSets []models.CalendarSet, workoutId uint) CalendarExerciseGroup {
	var exerciseGroup CalendarExerciseGroup
	if len(exercisesAtThisOrder) >= 2 {
		exerciseGroup = ConvertToCalendarSuperset(exercisesAtThisOrder, calendarSets, workoutId)
	} else {
		exerciseGroup = ConvertToCalendarSingleExercise(exercisesAtThisOrder, calendarSets, workoutId)
	}
	return exerciseGroup
}

func ConvertToCalendarSuperset(exercisesAtThisOrder []models.CalendarExercise, calendarSets []models.CalendarSet, workoutId uint) CalendarExerciseGroup {
	var exerciseGroup CalendarExerciseGroup
	orderedSupersetExercises := make([]models.CalendarExercise, len(exercisesAtThisOrder))
	for _, exercise := range exercisesAtThisOrder {
		orderedSupersetExercises[exercise.OrderInSuperset] = exercise
	}

	for _, exercise := range orderedSupersetExercises {
		relevantSets := GetRelevantSets(calendarSets, workoutId, exercise.ExerciseName)
		orderedRelevantSets := OrderRelevantSets(relevantSets)
		sets := ConvertSets(orderedRelevantSets)

		exerciseGroup = append(exerciseGroup, CalendarExercise{
			ExerciseName: exercise.ExerciseName,
			IsIsometric:  exercise.IsIsometric,
			Id:           exercise.ID,
			Sets:         sets,
		})
	}
	return exerciseGroup
}

func ConvertToCalendarSingleExercise(exercisesAtThisOrder []models.CalendarExercise, calendarSets []models.CalendarSet, workoutId uint) CalendarExerciseGroup {
	relevantSets := GetRelevantSets(calendarSets, workoutId, exercisesAtThisOrder[0].ExerciseName)
	orderedRelevantSets := OrderRelevantSets(relevantSets)
	sets := ConvertSets(orderedRelevantSets)

	return CalendarExerciseGroup{
		CalendarExercise{
			ExerciseName: exercisesAtThisOrder[0].ExerciseName,
			IsIsometric:  exercisesAtThisOrder[0].IsIsometric,
			Id:           exercisesAtThisOrder[0].ID,
			Sets:         sets,
		},
	}
}

func GetRelevantSets(calendarSets []models.CalendarSet, workoutId uint, exerciseName string) []models.CalendarSet {
	var relevantSets []models.CalendarSet
	for _, set := range calendarSets {
		if set.WorkoutId == workoutId && set.ExerciseName == exerciseName {
			relevantSets = append(relevantSets, set)
		}
	}
	return relevantSets
}

func OrderRelevantSets(calendarSets []models.CalendarSet) []models.CalendarSet {
	orderedRelevantSets := make([]models.CalendarSet, len(calendarSets))
	for _, set := range calendarSets {
		orderedRelevantSets[set.OrderInExercise] = set
	}
	return orderedRelevantSets
}

func ConvertSets(calendarSets []models.CalendarSet) Sets {
	var sets Sets
	for _, set := range calendarSets {
		sets = append(sets, Set{
			Weight:            set.Weight,
			KgOrLbs:           set.KgOrLbs,
			Reps:              set.Reps,
			IsometricHoldTime: set.IsometricHoldTime,
			Id:                set.ID,
		})
	}
	return sets
}
