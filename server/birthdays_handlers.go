package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/server/views"
)

// GET /birthdays
func (s *Server) birthdaysHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /birthdays")

	bdays, err := s.state.storage.GetAll(10, 0)
	if err != nil {
		slog.Error(
			"failed to get birthdays from db",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	birthdaysPage := views.Birthdays(bdays)

	w.Header().Add("Content-Type", "text/html")

	err = birthdaysPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Birthdays", slog.String("err_msg", err.Error()))
	}
}

// POST /api/birthdays
func (s *Server) handleCreateBirthday(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST /api/birthdays")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	day := r.FormValue("day")
	month := r.FormValue("month")

	d, err := strconv.ParseUint(day, 10, 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: validation
	bday := calendar.NewBirthday(name, time.Month(m), uint8(d))

	err = s.state.storage.Save(bday)
	if err != nil {
		slog.Error("failed to save bday", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GET /api/birthdays
// return todays birthdays (empty if no birthdays found)
func (s *Server) handleGetTodayBirthdays(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /api/birthdays")
	now := time.Now()
	day := now.Day()
	month := now.Month()

	birthdays, err := s.state.storage.GetAllForDate(month, uint8(day))
	if err != nil {
		slog.Error("failed to get otday birthdays", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calendarEvents := views.CalendarEvents(birthdays)
	w.Header().Add("Content-Type", "text/html")

	err = calendarEvents.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Birthdays", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
