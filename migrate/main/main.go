package main

import (
	"log"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
)

func init() {
	initializers.InitLogs()
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {
	//users
	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//templates
	err = initializers.DB.AutoMigrate(&models.TemplateSplit{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.TemplateSplitWorkoutLink{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.TemplateWorkout{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.TemplateExercise{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//calendar
	err = initializers.DB.AutoMigrate(&models.CalendarWorkout{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.CalendarExercise{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.CalendarSet{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//audits & rate limits
	err = initializers.DB.AutoMigrate(&models.Audit{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initializers.DB.AutoMigrate(&models.RateLimitAction{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//need user chart configs...
}
