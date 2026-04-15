package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leovaldes-debug/booking-service/internal/model"
)

// BookingRepository — слой работы с БД.
// Вся SQL-логика живёт здесь и нигде больше.
// Handler и Service про SQL не знают — это называется разделение ответственности.
type BookingRepository struct {
	db *pgxpool.Pool // pgxpool — пул соединений, не открываем новое соединение на каждый запрос
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create вставляет новое бронирование и сразу возвращает его с id и created_at.
// RETURNING — фишка PostgreSQL, не нужен отдельный SELECT после INSERT.
func (r *BookingRepository) Create(ctx context.Context, req model.CreateBookingRequest) (model.Booking, error) {
	var b model.Booking
	err := r.db.QueryRow(ctx,
		`INSERT INTO bookings (user_id, train_id, status)
		 VALUES ($1, $2, 'pending')
		 RETURNING id, user_id, train_id, status, created_at`,
		req.UserID, req.TrainID,
	).Scan(&b.ID, &b.UserID, &b.TrainID, &b.Status, &b.CreatedAt)
	if err != nil {
		return model.Booking{}, fmt.Errorf("create booking: %w", err)
	}
	return b, nil
}

// GetByID достаёт одно бронирование по id.
func (r *BookingRepository) GetByID(ctx context.Context, id int) (model.Booking, error) {
	var b model.Booking
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, train_id, status, created_at
		 FROM bookings WHERE id = $1`,
		id,
	).Scan(&b.ID, &b.UserID, &b.TrainID, &b.Status, &b.CreatedAt)
	if err != nil {
		return model.Booking{}, fmt.Errorf("get booking %d: %w", id, err)
	}
	return b, nil
}

// ListByUser возвращает все бронирования пользователя, новые сверху.
func (r *BookingRepository) ListByUser(ctx context.Context, userID int) ([]model.Booking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, train_id, status, created_at
		 FROM bookings WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.TrainID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// UpdateStatus меняет статус бронирования.
// Например: pending → confirmed, или pending → cancelled.
func (r *BookingRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bookings SET status = $1 WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update booking %d status: %w", id, err)
	}
	return nil
}
