package middleware

import (
	"github.com/gin-gonic/gin"
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
	initializers.DB.First(&user, userId)
	if user.ID == 0 || !user.IsVerified {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user.ID)
	c.Next()
}
