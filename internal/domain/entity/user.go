package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"column:id;type:char(36);primaryKey;not null"`
	Email        string         `gorm:"column:email;type:varchar(255);unique;not null"`
	Password     *string        `gorm:"column:password;type:varchar(255)"`
	Name         string         `gorm:"column:name;type:varchar(255);not null"`
	GoogleID     *string        `gorm:"column:google_id;type:varchar(255);unique"`
	Verified     bool           `gorm:"column:verified;type:bool;default:false"`
	Preference   []Preference   `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	RefreshToken []RefreshToken `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	CreatedAt    *time.Time     `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	u.ID = id
	return
}
