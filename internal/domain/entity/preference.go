package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Preference struct {
	ID        uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	Name      string     `gorm:"column:name;type:varchar(255);not null"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:char(36);not null"`
	User      *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (p *Preference) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	p.ID = id
	return
}
