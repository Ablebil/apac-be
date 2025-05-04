package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name       string    `gorm:"type:varchar(255);default:null"`
	Username   string    `gorm:"type:varchar(255);unique;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Password   string    `gorm:"type:varchar(255)"`
	IsVerified bool      `gorm:"type:boolean;default:false"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) {
	id, _ := uuid.NewV7()
	u.ID = id
	return
}
