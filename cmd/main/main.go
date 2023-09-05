package main

import (
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/services"
)

func init() {
	initializers.InitLogs()
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {
	go services.DeleteExpiredUnverifiedUsers()
	go services.RemoveExpiredPasswordResetCodes()
	InitRoutes()
}
