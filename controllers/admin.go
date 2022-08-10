package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

func GetAdminApproval(c *gin.Context) {

	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	username := authPayload.Username
	hackathonid := c.Params.ByName("hackathon_id")
	value := c.Params.ByName("value")
	user, err := repo.GetProfileByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	hackathon, err := repo.HackathonFromHackathonID(hackathonid)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	organiser, err := repo.UserFromUserID(hackathon.OrganiserID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.IsAdmin == true {
		err := repo.AdminAppoved(hackathonid, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Hackathon doesn't exist",
			})
			return
		} else if value == "0" {
			err := repo.DeleteHackathonByAdmin(hackathonid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			utils.Notify(organiser.Email, "Approval", "Unfortunately your hackathon has been disapproved.")
			c.JSON(http.StatusOK, gin.H{
				"message": "Hackathon disapproved and deleted from database",
			})
			return
		} else {
			utils.Notify(organiser.Email, "Approval", "Yay! your hackathon has been approved.")
			c.JSON(http.StatusOK, gin.H{
				"message": "Hackathon approved",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Sorry, you don't have the right to organise a hackathon. Please buy our premium version to unlock all features.",
		})
	}
}
