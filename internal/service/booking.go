package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/leovaldes-debug/booking-service/internal/model"
	"github.com/leovaldes-debug/booking-service/internal/repository"
)

// BookingService — бизнес-логика.
// Repository знает КАК сохранить, Service знает НУЖНО ЛИ сохранять.
// Именно здесь валидация, проверки, бизнес-правила.
type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create(ctx context.Context, req model.CreateBookingRequest) (model.Booking, error) {
	if req.UserID <= 0 {
		return model.Booking{}, errors.New("invalid user_id")
	}
	if req.TrainID <= 0 {
		return model.Booking{}, errors.New("invalid train_id")
	}
	return s.repo.Create(ctx, req)
}

func (s *BookingService) GetByID(ctx context.Context, id int) (model.Booking, error) {
	if id <= 0 {
		return model.Booking{}, errors.New("invalid id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *BookingService) ListByUser(ctx context.Context, userID int) ([]model.Booking, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	return s.repo.ListByUser(ctx, userID)
}

// Cancel — отдельный метод для отмены, а не просто UpdateStatus.
// Инкапсулирует логику: нельзя отменить уже отменённое.
func (s *BookingService) Cancel(ctx context.Context, id int) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("cancel booking: %w", err)
	}
	if b.Status == "cancelled" {
		return errors.New("booking already cancelled")
	}
	return s.repo.UpdateStatus(ctx, id, "cancelled")
}
