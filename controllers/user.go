package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
	"win-a-thon/models"
	"win-a-thon/repo"
	"win-a-thon/token"
	"win-a-thon/utils"
)

type createUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	Password     string `json:"password" binding:"required,min=6"`
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	LinkedIn     string `json:"linked_in"`
	GitHub       string `json:"github"`
	WebLink      string `json:"web_link"`
	Organisation string `json:"organisation"`
	IsAdmin      bool   `json:"is_admin"`
}

//CreateUser ... Create User
func CreateUser(ctx *gin.Context) {
	var req createUserRequest
	fmt.Println("accept: ", ctx.Param("Accept"))
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := models.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		LinkedIn:       req.LinkedIn,
		GitHub:         req.GitHub,
		WebLink:        req.WebLink,
		Organisation:   req.Organisation,
		IsAdmin:        false,
	}

	fmt.Println(arg)
	err = repo.CreateUser(&arg)

	if err != nil {
		if err.(*mysql.MySQLError).Number == utils.DuplicateRecordErrorCode {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "User Already exists"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "User registered successfully"})
}

func ShowUserHackathons(c *gin.Context) {
	var userHackathons []models.Participant
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
	err = repo.ShowUserHackathons(&userHackathons, user.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		var temp models.Hackathon
		type hackathonDetail struct {
			HackathonId  uint      `json:"hackathon_id"`
			Title        string    `json:"title"`
			StartingTime time.Time `json:"starting_time"`
			EndingTime   time.Time `json:"ending_time"`
			ResultTime   time.Time `json:"result_time"`
			Description  string    `json:"description"`
		}
		var hackathons []hackathonDetail
		for i := range userHackathons {
			hid := strconv.FormatUint(uint64(userHackathons[i].HackathonId), 10)
			err, _ := repo.ViewHackathonDetails(&temp, hid)
			if err == nil {
				var hd hackathonDetail
				hd.HackathonId = temp.ID
				hd.Title = temp.Title
				hd.StartingTime = temp.StartingTime
				hd.EndingTime = temp.EndingTime
				hd.ResultTime = temp.ResultTime
				hd.Description = temp.Description
				hackathons = append(hackathons, hd)
			}
		}
		var activeHackathons, pastHackathons, futureHackathons []hackathonDetail
		for i := range hackathons {
			if time.Now().After(hackathons[i].StartingTime) && time.Now().Before(hackathons[i].EndingTime) {
				activeHackathons = append(activeHackathons, hackathons[i])
			} else if time.Now().Before(hackathons[i].StartingTime) {
				futureHackathons = append(futureHackathons, hackathons[i])
			} else {
				pastHackathons = append(pastHackathons, hackathons[i])
			}
		}

		c.JSON(http.StatusOK, gin.H{"active hackathons": activeHackathons, "future hackathons": futureHackathons,
			"past hackathons": pastHackathons})

	}
}

type userProfileResponse struct {
	Username     string `json:"username" binding:"required,alphanum"`
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	LinkedIn     string `json:"linked_in"`
	GitHub       string `json:"github"`
	WebLink      string `json:"web_link"`
	Organisation string `json:"organisation"`
}

func NewUserProfileResponse(user models.User) userProfileResponse {
	return userProfileResponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		LinkedIn:     user.LinkedIn,
		GitHub:       user.GitHub,
		WebLink:      user.WebLink,
		Organisation: user.Organisation,
	}
}

func ViewUserProfile(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	authorisedUsername := authPayload.Username

	user, err := repo.GetProfileByUsername(authorisedUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
	}

	c.JSON(http.StatusOK, NewUserProfileResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string              `json:"access_token"`
	User        userProfileResponse `json:"user"`
}

func LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Incorrect credentials"})
		return
	}

	user, err := repo.GetProfileByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"msg": "user doesn't exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Incorrect Password"})
		return
	}

	tokenMaker, err := token.NewJWTMaker(utils.TokenSymmetricKey)
	duration, err := time.ParseDuration(utils.AccessTokenDuration)
	accessToken, err := tokenMaker.CreateToken(
		user.Username,
		duration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        NewUserProfileResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func DeleteUserProfile(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	authorisedUsername := authPayload.Username

	err := repo.DeleteUserByUsername(authorisedUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Profile deleted successfully"})
}

func CreateAdminUser(c *gin.Context) {

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := models.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		LinkedIn:       req.LinkedIn,
		GitHub:         req.GitHub,
		WebLink:        req.WebLink,
		Organisation:   req.Organisation,
		IsAdmin:        true,
	}

	fmt.Println(arg)
	err = repo.CreateUser(&arg)

	if err != nil {
		if err.(*mysql.MySQLError).Number == utils.DuplicateRecordErrorCode {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "User Already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Admin User registered successfully"})
}

func UpdateUserProfile(c *gin.Context) {
	authPayload := c.MustGet(utils.AuthorizationPayloadKey).(*token.Payload)
	authorisedUsername := authPayload.Username

	var req utils.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Update body format is incorrect"})
		return
	}

	user, err := repo.UpdateUserProfile(authorisedUsername, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
	}
	c.JSON(http.StatusOK, NewUserProfileResponse(user))
}
