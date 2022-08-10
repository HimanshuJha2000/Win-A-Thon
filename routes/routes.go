package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"win-a-thon/controllers"
	"win-a-thon/middlewares"
	"win-a-thon/token"
	"win-a-thon/utils"
)

func Setup() (*gin.Engine, error) {
	r := gin.Default()

	tokenMaker, err := token.NewJWTMaker(utils.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "I'm alive! Are you?"})
	})

	users := r.Group("/users/")
	{
		users.POST("signup", controllers.CreateUser)
		users.POST("login", controllers.LoginUser)
		users.POST("admin", controllers.CreateAdminUser)
	}

	protectedUsers := r.Group("/users").Use(middlewares.Authorization(tokenMaker))
	{
		protectedUsers.GET("", controllers.ViewUserProfile)
		protectedUsers.GET("hackathons", controllers.ShowUserHackathons)
		protectedUsers.DELETE("", controllers.DeleteUserProfile)
		protectedUsers.PATCH("", controllers.UpdateUserProfile)
	}

	protectedHackathons := r.Group("/hackathons").Use(middlewares.Authorization(tokenMaker))
	{
		protectedHackathons.POST("", controllers.CreateHackathon)
		protectedHackathons.GET("/:hackathon_id/submissions", controllers.GetSubmissions)
		protectedHackathons.GET("/:hackathon_id/participants", controllers.GetParticipants)
		protectedHackathons.PATCH("/:hackathon_id/submissions/:username/judge", controllers.JudgeSubmission)
		protectedHackathons.POST("/:hackathon_id/notify/:username", controllers.NotifyParticipant)
		protectedHackathons.POST("/:hackathon_id/notify", controllers.NotifyAll)
		protectedHackathons.GET("/organise", controllers.ListOrganisedHackathons)
		protectedHackathons.GET("/:hackathon_id/prize_distribution", controllers.PrizeDistribution)
		protectedHackathons.POST("/:hackathon_id/participate", controllers.Participate)
		protectedHackathons.PATCH("/:hackathon_id/submit", controllers.UpdateSubmission)
		protectedHackathons.PATCH("/:hackathon_id/approve/:value", controllers.GetAdminApproval)
		protectedHackathons.GET("/unapproved", controllers.ListUnapprovedHackathons)
		protectedHackathons.GET("/:hackathon_id/user/:username/submission", controllers.GetSubmissionOfParticipant)
	}

	hackathons := r.Group("/hackathons")
	{
		hackathons.GET("", controllers.ListHackathons)
		hackathons.GET("/:hackathon_id", controllers.ViewHackathonDetails)
		hackathons.GET("/:hackathon_id/leaderboard", controllers.GetLeaderboard)
	}

	return r, nil
}
