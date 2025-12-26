package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/maximekuhn/metego/internal/calendar"
	"github.com/maximekuhn/metego/internal/server/views"
)

// GET /birthdays
func (s *Server) birthdaysHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /birthdays")

	birthdays, err := s.state.bdaysStorage.GetAll(r.Context(), 1000, 0)
	if err != nil {
		slog.Error(
			"failed to get birthdays from db",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sort.Slice(birthdays, func(i, j int) bool {
		return birthdays[i].Date.Before(birthdays[j].Date)
	})

	birthdaysPage := views.Birthdays(birthdays)

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
	bday := calendar.NewBirthday(0, name, time.Month(m), uint8(d))

	err = s.state.bdaysStorage.Save(r.Context(), bday)
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

	birhtdays, err := s.state.bdaysStorage.GetAllForDate(r.Context(), month, uint8(day))
	if err != nil {
		slog.Error("failed to get today birthdays", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("birthdays", slog.Int("count", len(birhtdays)))

	// TODO: move this somewhere else
	apts, err := s.state.aptsStorage.GetAllForDate(r.Context(), uint8(day), month, uint(year))
	if err != nil {
		slog.Error("failed to get today appointments", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("appointments", slog.Int("count", len(apts)))

	// TODO: move this somewhere else
	nameday, found, err := s.state.namedaysStorage.GetForDate(r.Context(), month, uint8(day))
	if err != nil {
		slog.Error("failed to get today nameday", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("nameday", slog.Bool("found", found))

	calendarEvents := views.CalendarEvents(birhtdays, apts, nameday)
	w.Header().Add("Content-Type", "text/html")

	err = calendarEvents.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render calendar", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDeleteBirthday(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("DELETE %s", r.URL))

	bdayIDStr := r.PathValue("birthdayID")
	if bdayIDStr == "" {
		http.Error(w, "Missing birthdayID path parameter", http.StatusBadRequest)
		return
	}
	bdayID, err := strconv.Atoi(bdayIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	found, err := s.state.bdaysStorage.Delete(r.Context(), bdayID)
	if err != nil {
		slog.Error(
			"failed to delete birthday",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, fmt.Sprintf("Birthday %d not found", bdayID), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
