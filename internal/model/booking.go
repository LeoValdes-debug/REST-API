package model

import "time"

// Booking — бронирование тренировки.
// Это просто структура данных — Go не использует классы,
// данные и логика разделены намеренно.
type Booking struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainID   int       `json:"train_id"`
	Status    string    `json:"status"`    // pending / confirmed / cancelled
	CreatedAt time.Time `json:"created_at"`
}

// CreateBookingRequest — то что клиент присылает в POST /bookings.
// Отдельная структура чтобы не светить лишние поля наружу.
type CreateBookingRequest struct {
	UserID  int `json:"user_id"`
	TrainID int `json:"train_id"`
}
