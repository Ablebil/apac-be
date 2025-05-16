package usecase

import (
	"apac/internal/app/trip/repository"
	"apac/internal/domain/dto"
	res "apac/internal/infra/response"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"
)

type TripUsecaseItf interface {
	GetTripById(userId uuid.UUID, tripId uuid.UUID) (map[string]interface{}, *res.Err)
	GetAllTrips(userId uuid.UUID) ([]dto.TripSummaryResponse, *res.Err)
	Delete(userId uuid.UUID, tripId uuid.UUID) *res.Err
}

type TripUsecase struct {
	tripRepository repository.TripRepositoryItf
}

func NewTripUsecase(tripRepository repository.TripRepositoryItf) TripUsecaseItf {
	return &TripUsecase{
		tripRepository: tripRepository,
	}
}

func (uc *TripUsecase) GetTripById(userId uuid.UUID, tripId uuid.UUID) (map[string]interface{}, *res.Err) {
	trip, err := uc.tripRepository.FindById(userId, tripId)
	if err != nil {
		return nil, res.ErrInternalServer("Failed to find trip")
	}

	if trip == nil {
		return nil, res.ErrNotFound("Trip not found")
	}

	resp := trip.ParseDTOGet()
	return resp, nil
}

func (uc *TripUsecase) GetAllTrips(userId uuid.UUID) ([]dto.TripSummaryResponse, *res.Err) {
	trips, err := uc.tripRepository.FindAll(userId)
	if err != nil {
		return nil, res.ErrInternalServer("Failed to find trip")
	}

	if len(trips) == 0 {
		return nil, res.ErrNotFound("Trip not found")
	}

	var resps []dto.TripSummaryResponse
	for _, trip := range trips {
		var resp dto.TripSummaryResponse
		dtoResp := trip.ParseDTOGet()
		dtoResp["id"] = trip.ID
		mapstructure.Decode(dtoResp, &resp)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *TripUsecase) Delete(userId uuid.UUID, tripId uuid.UUID) *res.Err {
	err := uc.tripRepository.Delete(userId, tripId)
	if err != nil {
		return res.ErrInternalServer("Failed to delete trip")
	}
	return nil
}
