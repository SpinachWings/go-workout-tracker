package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	isValid, expiration, userId, err := services.ParseToken(tokenString)
	if err != nil || !isValid || time.Now().Unix() > expiration {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	result := initializers.DB.First(&user, userId)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding user for auth with ID: %d: %s", userId, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if user.ID == 0 || !user.IsVerified {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user.ID)
	c.Next()
}
