package models

type Participant struct {
	HackathonId int       `json:"hackathon_id" gorm:"primaryKey"`
	Hackathon   Hackathon `gorm:"foreignKey:HackathonId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId      int       `json:"user_id" gorm:"primaryKey"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DemoUrl     string    `json:"demo_url"`
	CodeUrl     string    `json:"code_url"`
	Score       int       `json:"score"`
}
