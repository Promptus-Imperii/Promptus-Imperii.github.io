package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/k42-software/go-altcha" // altcha
)

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
	// May fail if there the env variables are already loaded
	godotenv.Load()
	_, emailAddressExists := os.LookupEnv("EMAIL_ADDRESS")
	_, emailPasswordExists := os.LookupEnv("EMAIL_PASSWORD")
	if !(emailAddressExists && emailPasswordExists) {
		log.Fatalf("Emailadress and/or emailpassword environmentvariables not set")
	}

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("App has started, logging to file. Gin running in %s mode", gin.Mode())
	r := initRouter()

	c := cors.DefaultConfig()
	c.AllowAllOrigins = true

	r.Use(cors.New(c))

	r.Run(":3000")
}

func generateCaptchaChallenge(context *gin.Context) {
	challenge := altcha.NewChallengeEncoded()
	fmt.Println(challenge)
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

	email, _ := os.LookupEnv("EMAIL_ADDRESS")
	context.Data(http.StatusOK, "text/plain", []byte(email))
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

	fmt.Println(len(errors))
	if len(errors) != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"Errors": errors})
		return
	}

	exception_mail := os.Getenv("EMAIL_ADDRESS")
	memberErr := SendMemberInfoEmail(member)
	confirmationErr := SendNotificationEmail(member)

	if memberErr != nil || confirmationErr != nil {
		log.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"Errors": []string{fmt.Sprintf("Er is iets fout gegaan tijdens het verwerken van je aanmelden. Meld jezelf aan via %s", exception_mail)}})
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
