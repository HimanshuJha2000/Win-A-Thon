package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `json:"username" gorm:"type:varchar(50);not null;unique"`
	FullName       string `json:"full_name" gorm:"type:varchar(50)"`
	HashedPassword string `json:"hashed_password" gorm:"type:varchar(1000);not null"`
	Email          string `json:"email" gorm:"type:varchar(50);not null;unique"`
	LinkedIn       string `json:"linked_in"`
	GitHub         string `json:"git_hub"`
	WebLink        string `json:"web_link"`
	Organisation   string `json:"organisation"`
	IsAdmin        bool   `json:"is_admin"`
}
