package utils

import (
	"time"
)

func CurrentTimePlusMinutesAsUnix(minutes int) int64 {
	return time.Now().Add(time.Minute * time.Duration(minutes)).Unix()
}

func CurrentTimeMinusMinutesAsTime(minutes int) time.Time {
	return time.Now().Add(time.Minute * time.Duration(-minutes))
}

func CurrentTimeMinusHoursAsTime(minutes int) time.Time {
	return time.Now().Add(time.Hour * time.Duration(-minutes))
}

func SleepForHours(hours int) {
	time.Sleep(time.Hour * time.Duration(hours))
}

func DateAsStringIsInFuture(date string) bool {
	currentTime := time.Now()
	currentDate := currentTime.Format("2020-01-30")
	return date > currentDate
}

func DateAsStringIsMoreThanNumYearsInFuture(date string, years int) bool {
	futureTime := time.Now().AddDate(years, 0, 0)
	parsedDate, _ := time.Parse("2020-01-30", date)
	return parsedDate.After(futureTime)
}

func DateAsStringIsLessThanNumYearsInPast(date string, years int) bool {
	pastTime := time.Now().AddDate(-years, 0, 0)
	parsedDate, _ := time.Parse("2020-01-30", date)
	return pastTime.Before(parsedDate)
}
