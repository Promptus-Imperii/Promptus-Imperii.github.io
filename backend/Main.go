package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/k42-software/go-altcha" // altcha
)

var CORRESPONDANCE_EMAIL = ""
var SERVER_EMAIL_CREDENTIALS ServerEmailCredentials

func initRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	config := cors.DefaultConfig()

	if gin.Mode() == gin.DebugMode {
		config.AllowOrigins = []string{"*"}
	} else {
		config.AllowOrigins = []string{"https://beta.svpromptusimperii.nl", "https://svpromptusimperii.nl"}
	}

	router.Use(cors.New(config))
	api := router.Group("/api")
	api.GET("/captcha-challenge", generateCaptchaChallenge)
	api.POST("/signup", handleSignUp)
	api.POST("/email", getEmail)

	return router
}

func main() {
	godotenv.Load()

	// Fail early if the environment variables are not loaded
	serverEmail, serverEmailAddressExists := os.LookupEnv("SERVER_EMAIL_ADDRESS")
	emailPassword, emailPasswordExists := os.LookupEnv("EMAIL_PASSWORD")
	correspondanceEmail, correspondanceEmailAddressExists := os.LookupEnv("CORRESPONDANCE_EMAIL_ADDRESS")

	if !(correspondanceEmailAddressExists && emailPasswordExists && serverEmailAddressExists) {
		log.Fatalf("SERVER_EMAIL_ADDRESS, EMAIL_PASSWORD and/or CORRESPONDANCE_EMAIL_ADDRESS environmentvariables not set")
	}

	SERVER_EMAIL_CREDENTIALS = ServerEmailCredentials{
		email:    serverEmail,
		password: emailPassword,
	}
	CORRESPONDANCE_EMAIL = correspondanceEmail

	// Set logging to export to both a logfile and to stdout (the terminal)
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	stdoutAndFile := io.MultiWriter(f, os.Stdout)
	log.SetOutput(stdoutAndFile)

	log.Printf("App has started, logging to file and stdout. Gin running in %s mode", gin.Mode())
	r := initRouter()

	c := cors.DefaultConfig()
	c.AllowAllOrigins = true

	r.Use(cors.New(c))

	r.Run(":3000")
}

func generateCaptchaChallenge(context *gin.Context) {
	challenge := altcha.NewChallengeEncoded()

	log.Printf("Sending Altcha challange: %s", challenge)

	jsonData := []byte(challenge)
	context.Data(http.StatusOK, "application/json", jsonData)
}

func getEmail(context *gin.Context) {
	var req EmailRequest

	err := json.NewDecoder(context.Request.Body).Decode(&req)

	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !altchaGuard(context, req.Altcha) {
		return
	}

	context.Data(http.StatusOK, "text/plain", []byte(CORRESPONDANCE_EMAIL))
}

func handleSignUp(context *gin.Context) {
	var member PISignUp

	err := json.NewDecoder(context.Request.Body).Decode(&member)

	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !altchaGuard(context, member.Altcha) {
		return
	}

	var errors []string
	// oh boy i love validating
	member.PostalCode, err = validatePostalCode(member.PostalCode)
	errors = appendError(errors, err)
	errors = appendError(errors, validateDate(member.DateOfBirth))
	errors = appendError(errors, validatePhoneNumber(member.Phone, "Jouw telefoonnummer"))
	errors = appendError(errors, validateIBAN(member.IBAN))
	errors = appendError(errors, validatePhoneNumber(member.EmergencyContactPhoneNumber, "Het telefoonnummer van je noodcontact"))
	errors = appendError(errors, validateEmail(member.Email))
	errors = appendError(errors, validateCohortYear(member.CohortYear))

	log.Println(len(errors))
	if len(errors) != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"Errors": errors})
		return
	}

	memberErr := SendMemberInfoEmail(member, SERVER_EMAIL_CREDENTIALS, CORRESPONDANCE_EMAIL)
	confirmationErr := SendNotificationEmail(member, SERVER_EMAIL_CREDENTIALS, CORRESPONDANCE_EMAIL)

	if memberErr != nil || confirmationErr != nil {
		if memberErr != nil {
			log.Println(memberErr.Error())
		}
		if confirmationErr != nil {
			log.Println(confirmationErr.Error())
		}
		context.JSON(http.StatusInternalServerError, gin.H{"Errors": []string{fmt.Sprintf("Er is iets fout gegaan tijdens het verwerken van je aanmelden. Meld jezelf aan via %s", CORRESPONDANCE_EMAIL)}})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Success": "Registration successful."})
}

func altchaGuard(context *gin.Context, payload string) bool {
	valid := altcha.ValidateResponse(payload, false)

	if !valid && gin.Mode() != gin.TestMode {
		log.Println("Invalid Altcha payload", valid)
		context.JSON(http.StatusBadRequest, gin.H{"Errors": []string{"een geldige captcha is vereist. Probeer de pagina te herladen (je formuliervelden blijven bestaan)"}})
		return false
	}

	log.Println("Valid Altcha payload", valid, payload)
	return true
}

func appendError(errorList []string, err error) []string {
	if err != nil {
		errorList = append(errorList, err.Error())
	}
	return errorList
}
