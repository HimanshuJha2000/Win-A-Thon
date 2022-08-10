package repo

import (
	"win-a-thon/database"
	"win-a-thon/models"
)

func UserFromUserID(userID int) (models.User, error) {
	var user models.User
	err := database.DB.Model(&models.User{}).Where("ID = ?", userID).First(&user).Error
	return user, err
}

func HackathonFromHackathonID(hackathon_id string) (models.Hackathon, error) {
	var hackathon models.Hackathon
	err := database.DB.Where("ID = ?", hackathon_id).First(&hackathon).Error
	return hackathon, err
}

func JudgeSubmisson(participantInput *models.Participant, hackathon_id string, username string) error {
	var user models.User
	var participant models.Participant

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	user_id := int(user.ID)

	if err := database.DB.Where("hackathon_id = ? AND user_id = ?", hackathon_id, user_id).First(&participant).Error; err != nil {
		return err
	}

	participant.Score = participantInput.Score
	if err := database.DB.Save(&participant).Error; err != nil {
		return err
	}
	return nil
}

func GetLeaderboard(participants *[]models.Participant, hackathon_id string) (err error) {
	err = database.DB.Order("score desc").Where("hackathon_id = ?", hackathon_id).Find(participants).Error
	return err
}

func GetSubmissions(participants *[]models.Participant, hackathon_id string) (err error) {
	err = database.DB.Where("hackathon_id = ?", hackathon_id).Find(participants).Error
	return err
}

func GetWinners(winners *[]models.Participant, hackathon_id string) (err error) {
	err = database.DB.Where("hackathon_id = ?", hackathon_id).Order("score desc").Limit(3).Find(winners).Error
	return err
}
