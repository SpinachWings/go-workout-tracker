package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"log"
	"net/http"
	"workout-tracker-go-app/pkg/constants"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/pkg/models"
	"workout-tracker-go-app/pkg/services"
	"workout-tracker-go-app/pkg/utils"
)

type loginSignupBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type emailVerificationBody struct {
	Email            string `json:"email" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"`
}

func Signup(c *gin.Context) {
	var body loginSignupBody
	err := c.ShouldBindJSON(&body)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if !utils.IsValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
		return
	}

	if !utils.IsValidPassword(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient password - must be between 8-30 characters and contain a number, a lower case and a capital letter"})
		return
	}

	var alreadyPresentUser models.User
	result := initializers.DB.First(&alreadyPresentUser, "email = ?", body.Email)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error trying to detemine whether user exists with email: %s: %s", body.Email, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if alreadyPresentUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email already exists"})
		return
	}

	hash, err := models.EncryptPassword(body.Password, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	verificationCode := randstr.String(40)

	err = services.SendVerificationEmail(verificationCode, body.Email)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to send verification email to: %s: %s", body.Email, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	user, err := models.CreateUser(body.Email, hash, verificationCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	models.CreateAudit(constants.GetAuditTypes().UserCreation, user.ID, "")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user created with email: %s - a verification email has been sent if this email address exists. If you do not verify your email address within %d hours, the user will be deleted.", user.Email, constants.GetExpiryCheckTimes().UserWithUnverifiedEmail.ExpiryTimeInHours)})
}

func VerifyEmail(c *gin.Context) {
	var body emailVerificationBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = services.VerifyEmail(body.VerificationCode, body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code / user is already verified / user doesn't exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func Login(c *gin.Context) {
	var body loginSignupBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding user with email: %s: %s", body.Email, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	if models.RateLimitIsExceeded(constants.GetRateLimitActionTypes().Login, user.ID, "") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "rate limit exceeded"})
		return
	}

	err = models.ComparePassword(user.Password, body.Password)
	if err != nil && utils.IsMismatchedHashAndPassword(err) {
		models.CreateRateLimitRecord(constants.GetRateLimitActionTypes().Login, user.ID, "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		log.Print(fmt.Sprintf("Password comparison failed for user with ID: %d: %s", user.ID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not verified"})
		return
	}

	tokenDurationInMinutes := 60 // 10 on prod
	tokenString, err := services.CreateToken(user.ID, tokenDurationInMinutes)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to create token for user with ID: %d: %s", user.ID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	models.CreateAudit(constants.GetAuditTypes().UserLogin, user.ID, "")

	models.ClearOldRateLimitRecords(constants.GetRateLimitActionTypes().Login, user.ID, "")

	c.SetSameSite(http.SameSiteLaxMode)
	secondsInMinute := 60
	c.SetCookie("Authorization", tokenString, tokenDurationInMinutes*secondsInMinute, "", "", false, true) //secure = true on prod
	c.JSON(http.StatusOK, gin.H{"token": tokenString})                                                     // eventually, we don't want to return in this in the body & have it http only
}

func ValidateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user authorized"})
}

func RefreshToken(c *gin.Context) {
	userId, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenDurationInMinutes := 60 // 10 on prod
	tokenString, err := services.CreateToken(userId.(uint), tokenDurationInMinutes)
	if err != nil {
		log.Print(fmt.Sprintf("Failed to create token for user with ID: %s: %s", userId, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	secondsInMinute := 60
	c.SetCookie("Authorization", tokenString, tokenDurationInMinutes*secondsInMinute, "", "", false, true) //secure = true on prod
	c.JSON(http.StatusOK, gin.H{"token": tokenString})                                                     // eventually, we don't want to return in this in the body & have it http only
}

func DeleteUser(c *gin.Context) {
	var body loginSignupBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Print(fmt.Sprintf("Error finding user with email: %s: %s", body.Email, result.Error.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	if models.RateLimitIsExceeded(constants.GetRateLimitActionTypes().Login, user.ID, "") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "rate limit exceeded"})
		return
	}

	err = models.ComparePassword(user.Password, body.Password)
	if err != nil && utils.IsMismatchedHashAndPassword(err) {
		models.CreateRateLimitRecord(constants.GetRateLimitActionTypes().Delete, user.ID, "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		log.Print(fmt.Sprintf("Password comparison failed for user with ID: %d: %s", user.ID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected server error"})
		return
	}

	services.DeleteUser(user.ID)

	models.CreateAudit(constants.GetAuditTypes().UserDeletion, user.ID, "")

	models.ClearOldRateLimitRecords(constants.GetRateLimitActionTypes().Delete, user.ID, "")

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
