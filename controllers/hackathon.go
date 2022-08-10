package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"time"
	"win-a-thon/models"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

func ListHackathons(c *gin.Context) {
	var hackathons []models.Hackathon
	err := repo.ListHackathons(&hackathons, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Service unavailable")
		return
	}

	if len(hackathons) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, "No active or upcoming hackathons")
		return
	}

	type hackathonConcise struct {
		ID uint `json:"id"`
		Title        string    `json:"title"`
		StartingTime time.Time `json:"starting_time"`
		EndingTime   time.Time `json:"ending_time"`
		Participants int64     `json:"participants"`
	}
	var hackathonConciseList []hackathonConcise

	for _, val := range hackathons {
		count, err := repo.CountParticipants(strconv.Itoa(int(val.ID)))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		hackathonConciseList = append(hackathonConciseList, hackathonConcise{
			val.ID,
			val.Title,
			val.StartingTime,
			val.EndingTime,
			count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "successful",
		"hackathons": hackathonConciseList,
	})
}

func CreateHackathon(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	user, err := repo.GetProfileByUsername(authPayload.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
		return
	}

	var hackathon models.Hackathon

	err = c.BindJSON(&hackathon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "field names incorrect"})
		return
	}

	if hackathon.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title cannot be empty"})
		return
	}

	if hackathon.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Description cannot be empty"})
		return
	}

	if hackathon.EndingTime.Before(hackathon.StartingTime) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ending time should be after starting time"})
		return
	}

	if hackathon.ResultTime.Before(hackathon.EndingTime) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Result time should be after ending time"})
		return
	}

	if hackathon.EndingTime.Sub(hackathon.StartingTime) < time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hackathon should be at least an hour long"})
		return
	}

	if hackathon.ResultTime.Sub(hackathon.EndingTime) < time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Judging period should be at least an hour long"})
		return
	}

	hackathon.OrganiserID = int(user.ID)

	if user.IsAdmin {
		hackathon.AdminApproved = true
	} else {
		hackathon.AdminApproved = false
	}

	err = repo.CreateHackathon(&hackathon)
	if err != nil {
		err, ok := err.(*mysql.MySQLError)
		if ok {
			if err.Number == utils.DuplicateRecordErrorCode {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Hackathon Title already exists"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	if user.IsAdmin {
		c.JSON(http.StatusOK, gin.H{"status": "Request sent successfully. As you are an admin, hackathon has been automatically approved"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Request sent successfully. Kindly wait for admin approval."})
	}
}

func ViewHackathonDetails(c *gin.Context) {
	var hackathon models.Hackathon
	id := c.Params.ByName("hackathon_id")

	err, count := repo.ViewHackathonDetails(&hackathon, id)

	var val int
	val, err = strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)

	} else if hackathon.ID != uint(val) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This hackathon doesn't exist",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id": hackathon.ID,
			"title":             hackathon.Title,
			"starting_time":     hackathon.StartingTime,
			"ending_time":       hackathon.EndingTime,
			"description":       hackathon.Description,
			"result_time":       hackathon.ResultTime,
			"organisation_name": hackathon.OrganisationName,
			"participants":      count,
		})
	}
}

func Participate(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	username := authPayload.Username
	user, err := repo.GetProfileByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var participant models.Participant
	participant.Score = 0
	id := c.Params.ByName("hackathon_id")
	var isOrganiser bool
	isOrganiser, err = repo.IsOrganiser(user.ID, id)
	if isOrganiser == true {
		c.JSON(http.StatusOK, gin.H{"status": "Can't participate own hackathon"})
		return
	}

	var statusHackathon models.Hackathon
	err, _ = repo.ViewHackathonDetails(&statusHackathon, id)

	if time.Now().After(statusHackathon.EndingTime) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "hackathon ended already!"})
		return
	}

	participant.UserId = int(user.ID)
	participant.HackathonId, err = strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	err = repo.CreateParticipant(&participant)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		var email Email
		email.Subject = "Participation in hackathon" + statusHackathon.Title
		email.Message = "You have successfully participated in hackathon" + statusHackathon.Title
		if err := utils.Notify(user.Email, email.Subject, email.Message); err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(500, gin.H{
				"message": "server unavailable :(",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "participated successfully",
			"hackathon_id": statusHackathon.ID,
			"title":        statusHackathon.Title})
	}
}

func GetParticipants(c *gin.Context) {
	hackathonID := c.Params.ByName("hackathon_id")

	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	user, err := repo.GetProfileByUsername(authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
		return
	}

	hackathon, err := repo.HackathonFromHackathonID(hackathonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "this hackathon doesn't exist"})
		return
	}

	if int(user.ID) != hackathon.OrganiserID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "only the organiser can view all participants"})
		return
	}

	var participants []models.Participant

	err = repo.GetSubmissions(&participants, hackathonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
		return
	}

	count, err := repo.CountParticipants(hackathonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
		return
	}

	type ParticpantsConcise struct {
		ID uint `json:"id"`
		Username     string `json:"username"`
		FullName     string `json:"full_name"`
		Email        string `json:"email"`
		LinkedIn     string `json:"linked_in"`
		Github       string `json:"github"`
		WebLink      string `json:"web_link"`
		Organisation string `json:"organisation"`
	}
	var ParticipantsInfo []ParticpantsConcise

	for i := range participants {
		user, err := repo.UserFromUserID(participants[i].UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
			return
		}
		ParticipantsInfo = append(ParticipantsInfo, ParticpantsConcise{
			user.ID,
			user.Username,
			user.FullName,
			user.Email,
			user.LinkedIn,
			user.GitHub,
			user.WebLink,
			user.Organisation,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_participants": count,
		"participants":       ParticipantsInfo,
	})
}

func ListUnapprovedHackathons(c *gin.Context) {
	var hackathons []models.Hackathon

	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	user, err := repo.GetProfileByUsername(authPayload.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "the server encountered an issue"})
		return
	}
	if !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "this API is for admins only"})
		return
	}

	err = repo.ListHackathons(&hackathons, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "service unavailable")
		return
	}

	if len(hackathons) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, "no un-approved hackathons")
		return
	}

	type hackathonConcise struct {
		ID               uint      `json:"id"`
		Title            string    `json:"title"`
		StartingTime     time.Time `json:"starting_time"`
		EndingTime       time.Time `json:"ending_time"`
		ResultTime       time.Time `json:"result_time"`
		OrganisationName string    `json:"organisation_name"`
		Description      string    `json:"description"`
	}
	var hackathonConciseList []hackathonConcise

	for _, val := range hackathons {
		hackathonConciseList = append(hackathonConciseList, hackathonConcise{
			val.ID,
			val.Title,
			val.StartingTime,
			val.EndingTime,
			val.ResultTime,
			val.OrganisationName,
			val.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "successful",
		"hackathons": hackathonConciseList,
	})
}

func ListOrganisedHackathons(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	username := authPayload.Username
	user, err := repo.GetProfileByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Service unavailable",
		})
		return
	}

	id := user.ID // We get this from middleware

	var hackathons []models.Hackathon

	err = repo.ListOrganisedHackathons(&hackathons, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Service unavailable"})
		return
	}

	if len(hackathons) == 0 {
		c.JSON(http.StatusOK, "No Hackathons organised")
		return
	}

	type HackathonsWrapper struct {
		ID uint `json:"id"`
		Title            string    `json:"title"`
		StartingTime     time.Time `json:"starting_time"`
		EndingTime       time.Time `json:"ending_time"`
		ResultTime       time.Time `json:"result_time"`
		Description      string    `json:"description"`
		OrganisationName string    `json:"organisation_name"`
		AdminApproved    bool      `json:"admin_approved"`
		Participants     int64     `json:"participants"`
	}
	var HackathonsWrapperList []HackathonsWrapper

	for _, val := range hackathons {
		count, err := repo.CountParticipants(strconv.Itoa(int(val.ID)))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		HackathonsWrapperList = append(HackathonsWrapperList, HackathonsWrapper{
			val.ID,
			val.Title,
			val.StartingTime,
			val.EndingTime,
			val.ResultTime,
			val.Description,
			val.OrganisationName,
			val.AdminApproved,
			count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "successful",
		"organised_hackathons": HackathonsWrapperList,
	})
}