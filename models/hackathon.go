package models

import (
	"gorm.io/gorm"
	"time"
)

type Hackathon struct {
	gorm.Model
	Title            string       `json:"title" gorm:"type:varchar(50);not null;unique"`
	StartingTime     time.Time    `json:"starting_time"`
	EndingTime       time.Time    `json:"ending_time"`
	ResultTime       time.Time    `json:"result_time"`
	OrganisationName string       `json:"organisation_name" gorm:"type:varchar(50)"`
	OrganiserID      int          `json:"organiser_id" gorm:"not null"`
	User             User         `gorm:"foreignKey:OrganiserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description      string       `json:"description" gorm:"not null"`
	AdminApproved    bool `json:"admin_approved"`
}
