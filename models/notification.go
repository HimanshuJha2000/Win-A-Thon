package models

type Notification struct {
	UserId   int    `json:"id" gorm:"primaryKey"`
	User     User   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeviceId string `json:"device_id" gorm:"not null"`
}
