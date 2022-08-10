package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

func UpdateSubmission(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	username := authPayload.Username
	user, err := repo.GetProfileByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var p models.Participant

	id := c.Params.ByName("hackathon_id")

	IsPart, err := repo.IsParticipant(user.ID, id)
	if IsPart != true {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user has not participated",
			"error": err})
		return
	}

	var statusHackathon models.Hackathon
	err, _ = repo.ViewHackathonDetails(&statusHackathon, id)
	if time.Now().Before(statusHackathon.StartingTime) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "hackathon not started yet!"})
		return
	}
	if time.Now().After(statusHackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "hackathon ended already!"})
		return
	}
	p.HackathonId, err = strconv.Atoi(id)
	p.UserId = int(user.ID)

	type Submission struct {
		DemoUrl string `json:"demo_url"`
		CodeUrl string `json:"code_url"`
	}
	var temp Submission
	c.BindJSON(&temp)

	err = repo.GetSubmission(&p, p.HackathonId, user.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	p.DemoUrl = temp.DemoUrl
	p.CodeUrl = temp.CodeUrl

	database.DB.Save(&p)

	c.JSON(http.StatusOK, gin.H{
		"hackathon_id": statusHackathon.ID,
		"Title":        statusHackathon.Title,
		"code_url":     p.CodeUrl,
		"demo_url":     p.DemoUrl,
		"message":      "submitted successfully",
	})
	return
}
