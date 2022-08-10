package repo

import (
	"win-a-thon/database"
	"win-a-thon/models"
)

func AdminAppoved(hackathonID string, value string) (err error) {
	var Hack models.Hackathon
	err = database.DB.Where("id = ?", hackathonID).First(&Hack).Error
	if err != nil {
		return err
	}
	if value == "0" {
		Hack.AdminApproved = false
	} else {
		Hack.AdminApproved = true
	}
	database.DB.Save(&Hack)
	return nil
}

func DeleteHackathonByAdmin(hackathonID string) (err error) {
	var hackathon models.Hackathon
	if err := database.DB.Where("id = ?", hackathonID).Delete(&hackathon).Error; err != nil {
		return err
	}
	return nil
}
