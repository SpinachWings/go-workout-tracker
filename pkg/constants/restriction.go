package constants

type Restriction struct {
	ID                 int
	UnsubscribedAmount int
	SubscribedAmount   int
}

type Restrictions struct {
	TemplateSplitsPerUser               Restriction
	TemplateSplitMaxDuration            Restriction
	TemplateWorkoutsPerUser             Restriction
	TemplateMaxExerciseGroupsPerWorkout Restriction
	TemplateMaxExercisesPerGroup        Restriction
	CalendarWorkoutMinDateYears         Restriction
	CalendarWorkoutMaxDateYears         Restriction
	CalendarMaxExerciseGroupsPerWorkout Restriction
	CalendarMaxExercisesPerGroup        Restriction
	CalendarMaxSetsPerExercise          Restriction

	ChartsPerUser Restriction
}

func GetRestrictions() Restrictions {
	return Restrictions{
		TemplateSplitsPerUser:               Restriction{2, 5, 99},
		TemplateSplitMaxDuration:            Restriction{3, 14, 99},
		TemplateWorkoutsPerUser:             Restriction{1, 14, 99},
		TemplateMaxExerciseGroupsPerWorkout: Restriction{4, 20, 99},
		TemplateMaxExercisesPerGroup:        Restriction{5, 5, 99},
		CalendarWorkoutMinDateYears:         Restriction{6, 10, 10},
		CalendarWorkoutMaxDateYears:         Restriction{7, 10, 10},
		CalendarMaxExerciseGroupsPerWorkout: Restriction{8, 20, 99},
		CalendarMaxExercisesPerGroup:        Restriction{9, 5, 99},
		CalendarMaxSetsPerExercise:          Restriction{10, 10, 99},

		ChartsPerUser: Restriction{11, 3, 99}, //unsub amount will be 0 eventually
	}
}

func (r Restriction) GetRestrictionAmount(isSubscribed bool) int {
	if isSubscribed {
		return r.SubscribedAmount
	} else {
		return r.UnsubscribedAmount
	}
}
