package rest

import (
	"apac/internal/app/trip/usecase"
	res "apac/internal/infra/response"
	"apac/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TripHandler struct {
	TripUsecase usecase.TripUsecaseItf
}

func NewTripHandler(routerGroup fiber.Router, tripUsecase usecase.TripUsecaseItf, m middleware.MiddlewareItf) {
	TripHandler := TripHandler{
		TripUsecase: tripUsecase,
	}

	routerGroup = routerGroup.Group("/trips")
	routerGroup.Get("/:id", m.Authentication, TripHandler.GetTripById)
	routerGroup.Get("/", m.Authentication, TripHandler.GetAllTrips)
	routerGroup.Delete("/:id", m.Authentication, TripHandler.Delete)
}

func (h TripHandler) GetTripById(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userID").(uuid.UUID)
	tripId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.Error(ctx, res.ErrBadRequest("Invalid trip id"))
	}

	trip, errs := h.TripUsecase.GetTripById(userId, tripId)
	if errs != nil {
		return res.Error(ctx, errs)
	}

	return res.SuccessResponse(ctx, "Trip retrieved successfully", trip)
}

func (h TripHandler) GetAllTrips(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userID").(uuid.UUID)

	trips, err := h.TripUsecase.GetAllTrips(userId)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Trips retrieved successfully", trips)
}

func (h TripHandler) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userID").(uuid.UUID)
	tripId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.Error(ctx, res.ErrBadRequest(""))
	}

	return h.TripUsecase.Delete(userId, tripId)
}
