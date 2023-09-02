package constants

type RateLimitActionType struct {
	ID                int
	Action            string
	RateLimit         int
	DurationInMinutes int
}

type RateLimitActionTypes struct {
	Login  RateLimitActionType
	Delete RateLimitActionType
}

func GetRateLimitActionTypes() RateLimitActionTypes {
	return RateLimitActionTypes{
		Login:  RateLimitActionType{0, "Log in attempt with wrong password", 10, 10},
		Delete: RateLimitActionType{1, "Delete attempt with wrong password", 10, 10},
	}
}
