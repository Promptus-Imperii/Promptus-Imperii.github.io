package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/k42-software/go-altcha" // altcha
)

func initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/captcha-challenge", generateCaptchaChallenge)
	router.POST("/signup", handleSignUp)
	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file, exiting.")
	}

	r := initRouter()

	// FIXME bad CORS policy
	c := cors.DefaultConfig()
	c.AllowAllOrigins = true

	r.Use(cors.New(c))
	r.Run(":8080")
}

func generateCaptchaChallenge(context *gin.Context) {
	challenge := altcha.NewChallengeEncoded()
	fmt.Println(challenge)
	jsonData := []byte(challenge)
	context.Data(http.StatusOK, "application/json", jsonData)
}

func handleSignUp(context *gin.Context) {
	var signup PISignUp
	err := json.NewDecoder(context.Request.Body).Decode(&signup)
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Replay attack protection is off due to a bug.
	// https://github.com/k42-software/go-altcha/issues/1
	fmt.Print(gin.Mode())
	valid := altcha.ValidateResponse(signup.Altcha, false)
	if !valid && gin.Mode() != gin.TestMode {
		log.Println("Invalid Altcha payload", valid, signup.Altcha)
		context.JSON(http.StatusBadRequest, gin.H{"Bad Request": "Cannot signup without proper Altcha payload."})
		return
	}
	log.Println("Valid Altcha payload", valid, signup.Altcha)
	// normalize to save some time on regex :D
	signup.PostalCode = strings.ReplaceAll(signup.PostalCode, " ", "")

	// oh boy i love validating
	err = validatePostalCode(signup.PostalCode)
	if err != nil {
		returnErr(context, err)
		return
	}

	// validate own phone number
	err = validatePhoneNumber(signup.Phone)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validateIBAN(signup.IBAN)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validatePhoneNumber(signup.EmergencyContactPhoneNumber)
	if err != nil {
		returnErr(context, errors.New("noodcontact: "+err.Error()))
		return
	}

	err = validateEmail(signup.Email)
	if err != nil {
		returnErr(context, err)
		return
	}

	// at this point everything *should* be okay
	// sending the message already might be early

	context.JSON(http.StatusOK, gin.H{"Success": "Registration successful."})
}

func returnErr(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
}
