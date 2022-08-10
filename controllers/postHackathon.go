package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"win-a-thon/models"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

func JudgeSubmission(c *gin.Context) {
	var participant models.Participant

	hackathon_id := c.Param("hackathon_id")
	username := c.Param("username")

	c.BindJSON(&participant)

	//authorization
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)

	LoggedInUser, err1 := repo.GetProfileByUsername(authPayload.Username)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user doesn't exist"})
		return
	}

	hackathon, err2 := repo.HackathonFromHackathonID(hackathon_id)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "This hackathon id doesnt exist"})
		return
	}

	if hackathon.OrganiserID != int(LoggedInUser.ID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised to judge participant"})
		return
	}

	//judge can give score only after ending of hackathon
	if !time.Now().After(hackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Can't judge right now. Try again once the hackathon ends!!"})
		return
	}

	//judge can give score only after ending of hackathon
	if !time.Now().After(hackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Can't judge right now. Try again once the hackathon ends!!"})
		return
	}

	err := repo.JudgeSubmisson(&participant, hackathon_id, username)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "service unavailable",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "submission judged successfully!",
		})
	}
}

func GetSubmissions(c *gin.Context) {
	hackathon_id := c.Params.ByName("hackathon_id")

	var participants []models.Participant

	type submissions struct {
		HackathonID int `json:"hackathon_id"`
		Username string `json:"username"`
		CodeUrl  string `json:"code_url"`
		DemoUrl  string `json:"demo_url"`
	}

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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised to view submissions"})
		return
	}

	//the organiser can get submissions only after the hackathon ends
	if !time.Now().After(hackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Can't view submissions right now. Try again once the hackathon ends!!"})
		return
	}

	err := repo.GetSubmissions(&participants, hackathon_id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "request unsuccessful!",
		})
		return
	} else {
		obj := make([]submissions, 0)

		for i := 0; i < len(participants); i++ {
			user, err := repo.UserFromUserID(participants[i].UserId)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			temp := submissions{participants[i].HackathonId, user.Username, participants[i].CodeUrl, participants[i].DemoUrl}
			obj = append(obj, temp)
		}

		//responding back

		if err != nil {
			c.AbortWithStatus(500)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"total_submissions": len(obj),
				"submissions":       obj,
			})
		}
	}
}

func GetLeaderboard(c *gin.Context) {

	type results struct {
		UserName string `json:"user_name"`
		FullName string `json:"full_name"`
		Score    int    `json:"score"`
	}

	hackathon_id := c.Params.ByName("hackathon_id")
	var participants []models.Participant

	err := repo.GetLeaderboard(&participants, hackathon_id)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "service unavailable"})
		return

	} else {
		obj := make([]results, 0)

		for i := 0; i < len(participants); i++ {
			user, err := repo.UserFromUserID(participants[i].UserId)
			if err != nil {
				c.AbortWithStatus(500)
			}
			temp := results{user.Username, user.FullName, participants[i].Score}
			obj = append(obj, temp)
		}

		hackathon, err := repo.HackathonFromHackathonID(hackathon_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server unavailable"})
			return
		}

		if time.Now().After(hackathon.ResultTime) {
			c.JSON(http.StatusOK, gin.H{
				"leaderboard": obj,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "results not declared yet!!",
			})
		}
	}
}

func PrizeDistribution(c *gin.Context) {
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised! You are not the organiser/judge."})
		return
	}

	//time constraint check

	if !time.Now().After(hackathon.ResultTime) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Results aren't declared yet!"})
		return
	}

	var winners []models.Participant
	err := repo.GetWinners(&winners, hackathon_id)
	if err != nil {
		fmt.Print(err)
		return
	}

	var winners_user []models.User
	var winners_usernames []string

	for i := 0; i < len(winners); i++ {
		temp, err := repo.UserFromUserID(winners[i].UserId)
		if err != nil {
			fmt.Println(err)
		}

		winners_user = append(winners_user, temp)
		winners_usernames = append(winners_usernames, temp.Username)
	}

	message := "winners "

	for i := 0; i < len(winners_user); i++ {
		err := utils.Notify(winners_user[i].Email, "Prize distribution", "Congratulations! You’ve won the hackathon. Please follow this link to get your cash reward.")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server unavailable"})
			return
		} else {
			message = message + strconv.Itoa(i+1) + ", "
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"winners": winners_usernames,
		"status":  message + "notified successfully!!",
	})
}

func GetSubmissionOfParticipant(c *gin.Context) {

	hackathon_id := c.Params.ByName("hackathon_id")
	username := c.Params.ByName("username")

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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorised to view submission!"})
		return
	}

	//time constraint
	if !time.Now().After(hackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Can't view the submission right now. Try again once the hackathon ends!!"})
		return
	}

	user, err := repo.GetProfileByUsername(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Username doesn't exist"})
		return
	}

	participant, err := repo.ParticipantFromUserID(hackathon_id, strconv.Itoa(int(user.ID)))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "participant with username doesn't exist"})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"hackathon_id": participant.HackathonId,
			"user name": username,
			"code_url":  participant.CodeUrl,
			"demo_url":  participant.DemoUrl,
		})
	}
}
