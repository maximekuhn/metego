package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/maximekuhn/metego/internal/calendar"
	"github.com/maximekuhn/metego/internal/server/views"
)

// GET /birthdays
func (s *Server) birthdaysHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /birthdays")

	birhtdays, err := s.state.bdaysStorage.GetAll(1000, 0)
	if err != nil {
		slog.Error(
			"failed to get birthdays from db",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	birthdaysPage := views.Birthdays(birhtdays)

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

	err = s.state.bdaysStorage.Save(bday)
	if err != nil {
		// TODO: check error type
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
	year := now.Year()

	birhtdays, err := s.state.bdaysStorage.GetAllForDate(month, uint8(day))
	if err != nil {
		slog.Error("failed to get today birthdays", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("birthdays", slog.Int("count", len(birhtdays)))

	// TODO: move this somewhere else
	apts, err := s.state.aptsStorage.GetAllForDate(uint8(day), month, uint(year))
	if err != nil {
		slog.Error("failed to get today appointments", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("appointments", slog.Int("count", len(apts)))

	calendarEvents := views.CalendarEvents(birhtdays, apts)
	w.Header().Add("Content-Type", "text/html")

	err = calendarEvents.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render calendar", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
