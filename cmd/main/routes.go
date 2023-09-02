package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"workout-tracker-go-app/pkg/controllers"
	"workout-tracker-go-app/pkg/middleware"
)

func InitRoutes() {
	router := gin.Default()

	router.POST("/user/signup", controllers.Signup)
	router.POST("/user/verify/email", controllers.VerifyEmail)
	router.POST("/user/login", controllers.Login)
	router.POST("/user/validate", middleware.RequireAuth, controllers.ValidateUser)
	router.POST("/user/delete", middleware.RequireAuth, controllers.DeleteUser)
	router.POST("/user/refresh/token", middleware.RequireAuth, controllers.RefreshToken)

	router.POST("/password/reset/send", controllers.SendPasswordResetEmail)
	router.POST("/password/reset/confirm", controllers.ResetPassword)

	router.PUT("/template/workouts", middleware.RequireAuth, controllers.PutTemplateWorkouts)
	router.GET("/template/workouts", middleware.RequireAuth, controllers.GetTemplateWorkouts)

	router.PUT("/template/splits", middleware.RequireAuth, controllers.PutTemplateSplits)
	router.GET("/template/splits", middleware.RequireAuth, controllers.GetTemplateSplits)

	router.PUT("/calendar/workouts", middleware.RequireAuth, controllers.PutCalendarWorkouts)
	router.GET("/calendar/workouts", middleware.RequireAuth, controllers.GetCalendarWorkouts)

	// exercise names for auto complete dropdown list - get from calendar & template workouts with user id - return as array
	//router.GET("/exercise/names", middleware.RequireAuth, controllers.GetExerciseNames)

	// charts???

	err := router.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(fmt.Sprintf("Error running server on port: %s", os.Getenv("PORT")))
	}
}
