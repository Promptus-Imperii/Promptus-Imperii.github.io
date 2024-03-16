package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	valid := altcha.ValidateResponse(signup.Altcha, false)

	if !valid && gin.Mode() != gin.TestMode {
		log.Println("Invalid Altcha payload", valid)
		context.JSON(http.StatusBadRequest, gin.H{"Errors": []string{"een captcha is vereist"}})
		return
	}

	log.Println("Valid Altcha payload", valid, signup.Altcha)

	var errors []string

	// oh boy i love validating
	errors = appendError(errors, validatePostalCode(signup.PostalCode))
	errors = appendError(errors, validateDate(signup.DateOfBirth))
	errors = appendError(errors, validatePhoneNumber(signup.Phone, "Jouw telefoonnummer"))
	errors = appendError(errors, validateIBAN(signup.IBAN))
	errors = appendError(errors, validatePhoneNumber(signup.EmergencyContactPhoneNumber, "Het telefoonnummer van je noodcontact"))
	errors = appendError(errors, validateEmail(signup.Email))

	fmt.Println(len(errors))
	if len(errors) != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"Errors": errors})
		return
	}

	// at this point everything *should* be okay
	// sending the message already might be early

	context.JSON(http.StatusOK, gin.H{"Success": "Registration successful."})
}

func appendError(errorList []string, err error) []string {
	if err != nil {
		errorList = append(errorList, err.Error())
	}
	return errorList
}
