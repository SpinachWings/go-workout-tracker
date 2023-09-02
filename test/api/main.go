package main

import (
	"fmt"
	"os"
	"strings"
	"workout-tracker-go-app/pkg/initializers"
	"workout-tracker-go-app/test/api/utils"
)

func init() {
	initializers.LoadEnvVars()
}

type testConstants struct {
	appUrl                string
	contentType           string
	userEmail             string
	userPassword          string
	emailVerificationCode string
	passwordResetCode     string
	token                 string
	userPasswordUpdated   string
}

type loginSignupRequestBody struct {
	Email    string
	Password string
}

type emailVerificationBody struct {
	Email            string
	VerificationCode string
}

type sendPasswordResetEmailBody struct {
	Email string
}

type passwordResetBody struct {
	Email            string
	VerificationCode string
	NewPassword      string
}

func main() {
	var tests testConstants
	tests.appUrl = fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
	tests.contentType = "application/json"
	tests.userEmail = os.Getenv("API_TEST_EMAIL")
	tests.userPassword = "ProperPW123!"
	tests.userPasswordUpdated = "ProperPW321!"

	//tests.testCreateUserInvalidBody() // we only need to test invalid body with 1 endpoint provided all endpoints use ShouldBindJson & binding: required
	//tests.testCreateUserInvalidEmail()
	//tests.testCreateUserInvalidPassword()
	//tests.testCreateUserSuccess()
	//tests.testCreateUserAlreadyExists()

	// must re-assign this property to reach further tests
	tests.emailVerificationCode = "D1pEvCYW4KxOJLq6Bkxl"

	// must put c.JSON(http.StatusOK, gin.H{"token": tokenString}) at the end of func Login

	//tests.testLoginUnverifiedEmail()
	//tests.testVerifyEmailInvalid()
	//tests.testVerifyEmailSuccess()
	//tests.testVerifyEmailAlreadyVerified()

	//tests.testLoginIncorrectCredentials()
	//tests.testLoginIncorrectCredentialsRateLimit()
	//tests.testLoginSuccess() // must delete rate limits in db to reach login success

	// must re-assign this property to reach further tests
	tests.token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTExMTE1NjIsInN1YiI6MjF9.8lsp7KHrIoWu5HpPtBti9JNlr-_iBx0KJ1611I7sxtE"

	//tests.validateUserInvalid()
	//tests.validateUserSuccess()

	// must put c.JSON(http.StatusOK, gin.H{"token": tokenString}) at the end of func RefreshToken

	//tests.refreshTokenInvalid()
	//tests.refreshTokenSuccess()

	//tests.testSendResetPasswordCodeInvalidEmail()
	//tests.testSendResetPasswordCodeSuccess()

	// must re-assign this property to reach further tests
	tests.passwordResetCode = "HHhotdZs8jS1onm8oJfc"

	//tests.testResetPasswordInvalid()
	//tests.testResetPasswordSuccess()

	//tests.testDeleteUserIncorrectCredentials()
	//tests.testDeleteUserIncorrectCredentialsRateLimit()

	// run this on its own to delete the user
	//tests.testDeleteUserSuccess()
}

// SIGNUP

func (tc testConstants) testCreateUserInvalidBody() {
	var requestBodies []interface{}
	invalidRequestBody1 := struct{ Something string }{
		Something: "something",
	}
	invalidRequestBody2 := loginSignupRequestBody{
		Email: tc.userEmail,
	}
	invalidRequestBody3 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "",
	}
	requestBodies = append(requestBodies, invalidRequestBody1)
	requestBodies = append(requestBodies, invalidRequestBody2)
	requestBodies = append(requestBodies, invalidRequestBody3)

	endpoint := fmt.Sprintf("%s/user/signup", tc.appUrl)

	for _, requestBody := range requestBodies {
		response := utils.DoPost(requestBody, endpoint, tc.contentType)
		if response.ResponseBody["error"] != "Invalid request body" {
			fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
			panic("testCreateUserInvalidBody failed")
		}
	}
}

func (tc testConstants) testCreateUserInvalidEmail() {
	requestBody := loginSignupRequestBody{
		Email:    "thisisnotanemail",
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/signup", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["error"] != "Invalid email address" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testCreateUserInvalidEmail failed")
	}
}

func (tc testConstants) testCreateUserInvalidPassword() {
	var requestBodies []loginSignupRequestBody
	invalidRequestBody1 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "password1",
	}
	invalidRequestBody2 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "Password",
	}
	invalidRequestBody3 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "1Aa",
	}
	invalidRequestBody4 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "Password1Password1Password1Password1",
	}
	requestBodies = append(requestBodies, invalidRequestBody1)
	requestBodies = append(requestBodies, invalidRequestBody2)
	requestBodies = append(requestBodies, invalidRequestBody3)
	requestBodies = append(requestBodies, invalidRequestBody4)

	endpoint := fmt.Sprintf("%s/user/signup", tc.appUrl)

	for _, requestBody := range requestBodies {
		response := utils.DoPost(requestBody, endpoint, tc.contentType)
		if response.ResponseBody["error"] != "Insufficient password - must be between 8-30 characters and contain a number, a lower case and a capital letter" {
			fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
			panic("testCreateUserInvalidPassword failed")
		}
	}
}

func (tc testConstants) testCreateUserSuccess() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/signup", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if !strings.HasPrefix(response.ResponseBody["message"], "User created with email: ") {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testCreateUserSuccess failed")
	}
}

func (tc testConstants) testCreateUserAlreadyExists() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/signup", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)

	if response.ResponseBody["error"] != "User with this email already exists" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testCreateUserAlreadyExists failed")
	}
}

// VERIFY EMAIL

func (tc testConstants) testLoginUnverifiedEmail() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/login", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["error"] != "User is not verified" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testLoginUnverifiedEmail failed")
	}
}

func (tc testConstants) testVerifyEmailInvalid() {
	var requestBodies []emailVerificationBody
	invalidRequestBody1 := emailVerificationBody{
		Email:            "something",
		VerificationCode: tc.emailVerificationCode,
	}
	invalidRequestBody2 := emailVerificationBody{
		Email:            tc.userEmail,
		VerificationCode: "invalid",
	}
	requestBodies = append(requestBodies, invalidRequestBody1)
	requestBodies = append(requestBodies, invalidRequestBody2)

	endpoint := fmt.Sprintf("%s/user/verify/email", tc.appUrl)

	for _, requestBody := range requestBodies {
		response := utils.DoPost(requestBody, endpoint, tc.contentType)
		if response.ResponseBody["error"] != "Invalid verification code / user is already verified / user doesn't exists" {
			fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
			panic("testVerifyEmailInvalid failed")
		}
	}
}

func (tc testConstants) testVerifyEmailSuccess() {
	requestBody := emailVerificationBody{
		Email:            tc.userEmail,
		VerificationCode: tc.emailVerificationCode,
	}
	endpoint := fmt.Sprintf("%s/user/verify/email", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)

	if response.ResponseBody["message"] != "Email verified successfully" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testCreateUserAlreadyExists failed")
	}
}

func (tc testConstants) testVerifyEmailAlreadyVerified() {
	requestBody := emailVerificationBody{
		Email:            tc.userEmail,
		VerificationCode: tc.emailVerificationCode,
	}
	endpoint := fmt.Sprintf("%s/user/verify/email", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)

	if response.ResponseBody["error"] != "Invalid verification code / user is already verified / user doesn't exists" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testCreateUserAlreadyExists failed")
	}
}

// LOGIN

func (tc testConstants) testLoginIncorrectCredentials() {
	var requestBodies []loginSignupRequestBody
	invalidRequestBody1 := loginSignupRequestBody{
		Email:    "something",
		Password: tc.userPassword,
	}
	invalidRequestBody2 := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "invalid",
	}
	requestBodies = append(requestBodies, invalidRequestBody1)
	requestBodies = append(requestBodies, invalidRequestBody2)

	endpoint := fmt.Sprintf("%s/user/login", tc.appUrl)

	for _, requestBody := range requestBodies {
		response := utils.DoPost(requestBody, endpoint, tc.contentType)
		if response.ResponseBody["error"] != "Invalid email or password" {
			fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
			panic("testLoginIncorrectCredentials failed")
		}
	}
}

func (tc testConstants) testLoginIncorrectCredentialsRateLimit() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: "invalid",
	}
	endpoint := fmt.Sprintf("%s/user/login", tc.appUrl)

	for i := 0; i <= 10; i++ {
		utils.DoPost(requestBody, endpoint, tc.contentType)
	}

	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["error"] != "Rate limit exceeded" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testLoginIncorrectCredentialsRateLimit failed")
	}
}

func (tc testConstants) testLoginSuccess() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/login", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["token"] == "" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testLoginSuccess failed")
	}
	fmt.Println(response.ResponseBody["token"])
}

// AUTH

func (tc testConstants) validateUserInvalid() {
	endpoint := fmt.Sprintf("%s/user/validate", tc.appUrl)
	response := utils.DoPostWithCookieNoBodyInResponse(nil, endpoint, tc.contentType, "Authorization", "invalid")
	if response.StatusCode != 401 {
		fmt.Println(response.StatusCode)
		panic("validateUserInvalid failed")
	}
}

func (tc testConstants) validateUserSuccess() {
	endpoint := fmt.Sprintf("%s/user/validate", tc.appUrl)
	response := utils.DoPostWithCookieNoBodyInResponse(nil, endpoint, tc.contentType, "Authorization", tc.token)
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		panic("validateUserSuccess failed")
	}
}

// REFRESH TOKEN

func (tc testConstants) refreshTokenInvalid() {
	endpoint := fmt.Sprintf("%s/user/refresh/token", tc.appUrl)
	response := utils.DoPostWithCookieNoBodyInResponse(nil, endpoint, tc.contentType, "Authorization", "invalid")
	if response.StatusCode != 401 {
		fmt.Println(response.StatusCode)
		panic("refreshTokenInvalid failed")
	}
}

func (tc testConstants) refreshTokenSuccess() {
	endpoint := fmt.Sprintf("%s/user/refresh/token", tc.appUrl)
	response := utils.DoPostWithCookie(nil, endpoint, tc.contentType, "Authorization", tc.token)
	if response.ResponseBody["token"] == "" {
		fmt.Println(response.StatusCode)
		panic("refreshTokenSuccess failed")
	}
	fmt.Println(response.ResponseBody["token"])
}

// RESET PASSWORD

func (tc testConstants) testSendResetPasswordCodeInvalidEmail() {
	requestBody := sendPasswordResetEmailBody{
		Email: "invalid",
	}
	endpoint := fmt.Sprintf("%s/password/reset/send", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["message"] != "Password reset email has been sent if user with this email exists" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testSendResetPasswordCodeInvalidEmail failed")
	}
}

func (tc testConstants) testSendResetPasswordCodeSuccess() {
	requestBody := sendPasswordResetEmailBody{
		Email: tc.userEmail,
	}
	endpoint := fmt.Sprintf("%s/password/reset/send", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["message"] != "Password reset email has been sent if user with this email exists" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testSendResetPasswordCodeSuccess failed")
	}
}

func (tc testConstants) testResetPasswordInvalid() {
	var requestBodies []passwordResetBody
	invalidRequestBody1 := passwordResetBody{
		Email:            "something",
		VerificationCode: tc.passwordResetCode,
		NewPassword:      tc.userPasswordUpdated,
	}
	invalidRequestBody2 := passwordResetBody{
		Email:            tc.userEmail,
		VerificationCode: "invalid",
		NewPassword:      tc.userPasswordUpdated,
	}
	requestBodies = append(requestBodies, invalidRequestBody1)
	requestBodies = append(requestBodies, invalidRequestBody2)

	endpoint := fmt.Sprintf("%s/password/reset/confirm", tc.appUrl)

	for _, requestBody := range requestBodies {
		response := utils.DoPost(requestBody, endpoint, tc.contentType)
		if response.ResponseBody["error"] != "Invalid verification code or email" {
			fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
			panic("testResetPasswordInvalid failed")
		}
	}
}

func (tc testConstants) testResetPasswordSuccess() {
	requestBody := passwordResetBody{
		Email:            tc.userEmail,
		VerificationCode: tc.passwordResetCode,
		NewPassword:      tc.userPasswordUpdated,
	}
	endpoint := fmt.Sprintf("%s/password/reset/confirm", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["message"] != "Password reset successful" {
		fmt.Println(response.ResponseBody["message"] + response.ResponseBody["error"])
		panic("testResetPasswordSuccess failed")
	}
}

// DELETE USER

func (tc testConstants) testDeleteUserSuccess() {
	requestBody := loginSignupRequestBody{
		Email:    tc.userEmail,
		Password: tc.userPassword,
	}
	endpoint := fmt.Sprintf("%s/user/delete", tc.appUrl)
	response := utils.DoPost(requestBody, endpoint, tc.contentType)
	if response.ResponseBody["error"] != "User deleted" {
		panic("testDeleteUser failed")
	}
}
