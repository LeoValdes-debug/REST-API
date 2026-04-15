package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/leovaldes-debug/booking-service/internal/handler"
	"github.com/leovaldes-debug/booking-service/internal/repository"
	"github.com/leovaldes-debug/booking-service/internal/service"
)

func main() {
	// Загружаем .env если он есть (в продакшне переменные придут из окружения)
	_ = godotenv.Load()

	// Подключаемся к PostgreSQL через пул соединений.
	// pgxpool — не открывает новое соединение на каждый запрос,
	// держит пул готовых соединений. Это важно под нагрузкой.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("connected to database")

	// Собираем зависимости вручную — dependency injection без фреймворков.
	// Repository → Service → Handler. Каждый слой знает только о следующем.
	repo := repository.NewBookingRepository(pool)
	svc := service.NewBookingService(repo)
	h := handler.NewBookingHandler(svc)

	// chi-роутер с базовыми middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // логирует каждый запрос
	r.Use(middleware.Recoverer) // ловит panic и возвращает 500 вместо краша

	r.Mount("/bookings", h.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown — ждём завершения текущих запросов перед остановкой.
	// Без этого при деплое можно оборвать запрос на полуслове.
	go func() {
		log.Printf("server started on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}
