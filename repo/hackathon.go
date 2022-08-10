package repo

import (
	"win-a-thon/database"
	"win-a-thon/models"
)

func ListHackathons(hackathon *[]models.Hackathon, adminApproved bool) error {
	err := database.DB.Where("admin_approved = ?", adminApproved).Order("starting_" +
		"time desc").Find(hackathon).Error
	return err
}

func ViewHackathonDetails(hackathon *models.Hackathon, HackathonID string) (error, int64) {
	if err := database.DB.Where("id = ? AND admin_approved = ?", HackathonID, true).First(hackathon).Error; err != nil {
		return err, 0
	}
	count, err := CountParticipants(HackathonID)
	if err != nil {
		return err, 0
	}
	return nil, count
}

func CreateHackathon(hackathon *models.Hackathon) error {
	err := database.DB.Create(hackathon).Error
	return err
}

func ListOrganisedHackathons(hackathon *[]models.Hackathon, id uint) error {
	err := database.DB.Where("organiser_id = ?", id).Order("starting_" +
		"time desc").Find(hackathon).Error
	return err
}

func CreateParticipant(participant *models.Participant) (err error) {
	result := database.DB.Create(participant)
	return result.Error
}
