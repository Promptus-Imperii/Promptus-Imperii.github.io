package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/k42-software/go-altcha"
	altchaHttp "github.com/k42-software/go-altcha/http" // altcha
)

func initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/altcha.min.js", func(c *gin.Context) {
		altchaHttp.ServeJavascript(c.Writer, c.Request)
	})
	router.GET("/captcha-challenge", generateCaptchaChallenge)
	router.POST("/signup", handleSignUp)
	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file, exiting.")
	}

	mondayToken = os.Getenv("MONDAYTOKEN")

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
	var signup PISignUP
	err := json.NewDecoder(context.Request.Body).Decode(&signup)
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := altcha.ValidateResponse(signup.Altcha, true)
	if !valid {
		log.Println("Invalid Altcha payload", valid, signup.Altcha)
		context.JSON(http.StatusBadRequest, gin.H{"Bad Request": "Cannot signup without proper Altcha payload."})
	}

	log.Println("Post data: ", signup)

	// normalize to save some time on regex :D
	signup.PostalCode = strings.ReplaceAll(signup.PostalCode, " ", "")

	// oh boy i love validating
	err = validatePostalCode(signup.PostalCode)
	if err != nil {
		returnErr(context, err)
		return
	}

	// validate own phone number
	err = validatePhoneNumber(signup.Member.PhoneNumber)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validateIBAN(signup.IBAN)
	if err != nil {
		returnErr(context, err)
		return
	}

	err = validatePhoneNumber(signup.EmergencyContact.PhoneNumber)
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
