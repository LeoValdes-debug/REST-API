-- Миграция создаёт таблицу бронирований.
-- Запускать: psql $DATABASE_URL -f migrations/001_create_bookings.sql

CREATE TABLE IF NOT EXISTS bookings (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    train_id   INT NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending'
                   CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс по user_id — запрос ListByUser будет частым
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
