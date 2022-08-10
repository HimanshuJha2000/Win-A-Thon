package repo

import (
	"fmt"
	"win-a-thon/database"
	"win-a-thon/models"
)

func CountParticipants(hackathonID string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Participant{}).Where("hackathon_id = ?", hackathonID).Count(&count).Error
	fmt.Println(count)
	return count, err
}

func ParticipantFromUserID(hackathon_id string, user_id string) (models.Participant, error) {
	var participant models.Participant
	if err := database.DB.Where("hackathon_id = ? AND user_id = ?", hackathon_id, user_id).Find(&participant).Error; err != nil {
		return participant, err
	}
	return participant, nil
}

func IsParticipant(uid uint, hid string) (bool, error) {
	var temp models.Participant
	if err := database.DB.Where("hackathon_id = ? AND user_id = ?", hid, uid).First(&temp).Error; err != nil {
		return false, err
	}
	return true, nil
}

func GetSubmission(p *models.Participant, id int, uid uint) (err error) {
	result := database.DB.Where("hackathon_id = ?", id).Where("user_id = ?", uid).First(&p)
	return result.Error
}

func IsOrganiser(uid uint, hid string) (bool, error) {
	var temp models.Hackathon
	result := database.DB.Where("id = ?", hid).Find(&temp)
	if uint(temp.OrganiserID) == uid {
		return true, result.Error
	}
	return false, result.Error
}
