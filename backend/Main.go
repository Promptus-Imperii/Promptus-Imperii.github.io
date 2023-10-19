package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/signup", handleSignUp)
	return router
}

func main() {
	r := initRouter()
	r.Run(":8080")
}

func handleSignUp(context *gin.Context) {
	var signup PISignUP

	err := json.NewDecoder(context.Request.Body).Decode(&signup)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO(jonas): make sure that if a validation fails,
	// a descriptive error is returned so a clear error message can be
	// displayed on the page.

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
