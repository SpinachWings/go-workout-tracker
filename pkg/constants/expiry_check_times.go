package constants

type ExpiryCheckTime struct {
	ID                int
	SleepTimeInHours  int
	ExpiryTimeInHours int
}

type ExpiryCheckTimes struct {
	UserWithUnverifiedEmail ExpiryCheckTime
	UserPasswordResetCode   ExpiryCheckTime
}

func GetExpiryCheckTimes() ExpiryCheckTimes {
	return ExpiryCheckTimes{
		UserWithUnverifiedEmail: ExpiryCheckTime{0, 24, 4},
		UserPasswordResetCode:   ExpiryCheckTime{1, 24, 4},
	}
}
