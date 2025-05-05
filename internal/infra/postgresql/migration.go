package postgresql

import (
	"apac/internal/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(entity.User{}, entity.RefreshToken{}, entity.Preference{})
}
