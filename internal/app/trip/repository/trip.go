package repository

import (
	"apac/internal/domain/entity"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TripRepositoryItf interface {
	Create(trip *entity.Trip) (*entity.Trip, error)
	FindById(userId uuid.UUID, tripId uuid.UUID) (*entity.Trip, error)
	FindAll(userId uuid.UUID) ([]entity.Trip, error)
	Delete(userId uuid.UUID, tripId uuid.UUID) error
}

type TripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepositoryItf {
	return &TripRepository{db}
}

func (t *TripRepository) Create(trip *entity.Trip) (*entity.Trip, error) {
	err := t.db.Clauses(clause.Returning{}).Select("Email", "Password").Create(trip).Error
	if err != nil {
		return nil, err
	}

	return trip, nil
}

func (t *TripRepository) FindById(userId uuid.UUID, tripId uuid.UUID) (*entity.Trip, error) {
	var trip entity.Trip
	err := t.db.Where("user_id = ?", userId).Where("id = ?", tripId).First(&trip).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &trip, nil
}

func (t *TripRepository) FindAll(userId uuid.UUID) ([]entity.Trip, error) {
	var trips []entity.Trip

	err := t.db.Where("user_id = ?", userId).Find(&trips).Error
	if err != nil {
		return nil, err
	}

	return trips, nil
}

func (t *TripRepository) Delete(userId uuid.UUID, tripId uuid.UUID) error {
	err := t.db.Where("user_id = ?", userId).Delete("id = ?", tripId).Error
	if err != nil {
		return err
	}

	return nil
}
