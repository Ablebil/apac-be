package entity

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Trip struct {
	ID      uuid.UUID `gorm:"column:id;type:char(36);primaryKey;not null"`
	UserID  uuid.UUID `gorm:"column:user_id;type:char(36);not null"`
	Content string    `gorm:"column:content;type:jsonb;not null"`
	User    *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (t *Trip) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	t.ID = id
	return
}

func (t *Trip) ParseDTOGet() map[string]interface{} {
	tripResponse := make(map[string]interface{})

	json.Unmarshal([]byte(t.Content), &tripResponse)

	return tripResponse
}
