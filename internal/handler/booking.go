package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/leovaldes-debug/booking-service/internal/model"
	"github.com/leovaldes-debug/booking-service/internal/service"
)

// BookingHandler — HTTP-слой.
// Его работа: распарсить запрос, вызвать сервис, вернуть ответ.
// Никакой бизнес-логики здесь нет — только HTTP.
type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// Routes регистрирует все маршруты на роутере chi.
// chi выбран потому что он легковесный, совместим со стандартной библиотекой
// и хорошо поддерживает middleware.
func (h *BookingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/user/{userID}", h.ListByUser)
	r.Delete("/{id}", h.Cancel)
	return r
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	booking, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	booking, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	bookings, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Возвращаем пустой массив вместо null — фронтенд скажет спасибо
	if bookings == nil {
		bookings = []model.Booking{}
	}

	writeJSON(w, http.StatusOK, bookings)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Cancel(r.Context(), id); err != nil {
		if errors.Is(err, errors.New("booking already cancelled")) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON — хелпер чтобы не повторять json.Marshal везде
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
