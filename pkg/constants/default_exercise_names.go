package constants

type DefaultExerciseName struct {
	ID          int
	Description string
}

type DefaultExerciseNames struct {
	BarbellBenchPress       DefaultExerciseName
	DumbbellBenchPress      DefaultExerciseName
	InclineBarbellPress     DefaultExerciseName
	InclineDumbbellPress    DefaultExerciseName
	DumbbellFlye            DefaultExerciseName
	InclineDumbbellFlye     DefaultExerciseName
	PushUp                  DefaultExerciseName
	InclinePushUp           DefaultExerciseName
	RingPushUp              DefaultExerciseName
	InclineRingPushUp       DefaultExerciseName
	ParallettePushUp        DefaultExerciseName
	InclineParallettePushUp DefaultExerciseName
}

func GetDefaultExerciseNames() DefaultExerciseNames {
	return DefaultExerciseNames{
		BarbellBenchPress:  DefaultExerciseName{0, "Barbell bench press"},
		DumbbellBenchPress: DefaultExerciseName{1, "Dumbbell bench press"},

		// etc... be sure to output these with exercise name controller
	}
}
