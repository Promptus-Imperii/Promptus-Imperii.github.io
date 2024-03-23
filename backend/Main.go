package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/k42-software/go-altcha" // altcha
)

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/captcha-challenge", generateCaptchaChallenge)
	api.POST("/signup", handleSignUp)
	return router
}

func main() {
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
	var member PISignUp
	err := json.NewDecoder(context.Request.Body).Decode(&member)
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Replay attack protection is off due to a bug.
	// https://github.com/k42-software/go-altcha/issues/1
	valid := altcha.ValidateResponse(member.Altcha, false)

	if !valid && gin.Mode() != gin.TestMode {
		log.Println("Invalid Altcha payload", valid)
		context.JSON(http.StatusBadRequest, gin.H{"Errors": []string{"een geldige captcha is vereist. Probeer de pagina te herladen (je formuliervelden blijven bestaan)"}})
		return
	}

	log.Println("Valid Altcha payload", valid, member.Altcha)

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

	// at this point everything *should* be okay
	// sending the message already might be early
	if gin.Mode() != gin.TestMode {
		SendMember(member)
	}

	context.JSON(http.StatusOK, gin.H{"Success": "Registration successful."})
}

func appendError(errorList []string, err error) []string {
	if err != nil {
		errorList = append(errorList, err.Error())
	}
	return errorList
}
