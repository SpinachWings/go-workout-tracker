package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"workout-tracker-go-app/pkg/controllers"
	"workout-tracker-go-app/pkg/middleware"
)

func InitRoutes() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:           []string{os.Getenv("CLIENT_ORIGIN")},
		AllowMethods:           []string{"*"},
		AllowHeaders:           []string{"Access-Control-Allow-Headers", "Origin,Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             true,
	}))

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

	router.GET("/exercise/names", middleware.RequireAuth, controllers.GetExerciseNames)

	// charts???

	err := router.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(fmt.Sprintf("Error running server on port: %s", os.Getenv("PORT")))
	}
}
