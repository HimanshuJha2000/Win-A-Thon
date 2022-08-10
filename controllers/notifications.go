package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"win-a-thon/models"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

type Email struct {
	Subject string
	Message string
}

func NotifyParticipant(c *gin.Context) {
	var email Email

	c.BindJSON(&email)

	username := c.Params.ByName("username")
	hackathon_id := c.Params.ByName("hackathon_id")

	//authorization
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)

	LoggedInUser, err1 := repo.GetProfileByUsername(authPayload.Username)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User doesn't exist"})
		return
	}

	hackathon, err2 := repo.HackathonFromHackathonID(hackathon_id)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "This hackathon doesnt exist"})
		return
	}

	if hackathon.OrganiserID != int(LoggedInUser.ID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised to notify participant"})
		return
	}

	//getting user from user name
	user, err1 := repo.GetProfileByUsername(username)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "This username doesn't exist"})
		return
	}

	//calling the notify function using the details of the user and responding back
	if err := utils.Notify(user.Email, email.Subject, email.Message); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "server unavailable :(",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "participant notified successfully!",
		})
	}
}

func NotifyAll(c *gin.Context) {
	var email Email
	var participants []models.Participant

	c.BindJSON(&email)

	hackathon_id := c.Params.ByName("hackathon_id")

	//authorization
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)

	LoggedInUser, err1 := repo.GetProfileByUsername(authPayload.Username)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User doesn't exist"})
		return
	}

	hackathon, err2 := repo.HackathonFromHackathonID(hackathon_id)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "This hackathon id doesnt exist"})
		return
	}

	if hackathon.OrganiserID != int(LoggedInUser.ID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised to notify participant"})
		return
	}

	//retrieving a slice of participants
	err := repo.GetSubmissions(&participants, hackathon_id)

	var err3 error
	for i := 0; i < len(participants); i++ {

		//retrieving user from user id
		var user models.User
		user, err := repo.UserFromUserID(participants[i].UserId)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		//call notify function to send the emails
		if err3 = utils.Notify(user.Email, email.Subject, email.Message); err3 != nil {
			break
		}
	}

	//responding back
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "server unavailable :(",
		})
	} else if err1 != nil {
		fmt.Println(err1.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "server unavailable :(",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "all participants notified successfully!",
		})
	}
}
