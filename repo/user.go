package repo

import (
	"win-a-thon/database"
	"win-a-thon/models"
	"win-a-thon/utils"
)

func CreateUser(user *models.User) (err error) {
	if err = database.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func ShowUserHackathons(userHackathons *[]models.Participant, uid uint) (err error) {
	result := database.DB.Where("user_id = ?", uid).Find(userHackathons)
	return result.Error
}

func GetProfileByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func DeleteUserByUsername(username string) error {
	var user models.User
	result := database.DB.Where("username = ?", username).Delete(&user)
	return result.Error
}

func UpdateUserProfile(username string, req utils.UpdateProfileRequest) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "username = ?", username)

	updatedEmail := user.Email
	updatedLinkedIn := user.LinkedIn
	updatedGitHub := user.GitHub
	updatedWebLink := user.WebLink
	updatedOrganisation := user.Organisation
	if req.Email != "" {
		updatedEmail = req.Email
	}
	if req.LinkedIn != "" {
		updatedLinkedIn = req.LinkedIn
	}
	if req.GitHub != "" {
		updatedGitHub = req.GitHub
	}
	if req.WebLink != "" {
		updatedWebLink = req.WebLink
	}
	if req.Organisation != "" {
		updatedOrganisation = req.Organisation
	}

	database.DB.Model(&user).Updates(models.User{Email: updatedEmail, LinkedIn: updatedLinkedIn, GitHub: updatedGitHub, WebLink: updatedWebLink, Organisation: updatedOrganisation})

	return user, result.Error
}
