package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/maximekuhn/metego/internal/calendar"
	"github.com/maximekuhn/metego/internal/server/views"
)

// GET /appointments
func (s *Server) appointmentsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /appointments")

	apts, err := s.state.aptsStorage.GetAll(255, 0)
	if err != nil {
		slog.Error(
			"failed to get appointments from db",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	aptsPage := views.Appointments(apts)

	w.Header().Add("Content-Type", "text/html")

	err = aptsPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Appointments", slog.String("err_msg", err.Error()))
	}
}

func (s *Server) handleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST /api/appointments")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	day := r.FormValue("day")
	month := r.FormValue("month")
	year := r.FormValue("year")
	hour := r.FormValue("hour")
	minute := r.FormValue("minute")

	// try to create a time.Time from the input
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		http.Error(w, "Invalid day value", http.StatusBadRequest)
		return
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		http.Error(w, "Invalid month value", http.StatusBadRequest)
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "Invalid year value", http.StatusBadRequest)
		return
	}
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		http.Error(w, "Invalid hour value", http.StatusBadRequest)
		return
	}
	minuteInt, err := strconv.Atoi(minute)
	if err != nil {
		http.Error(w, "Invalid minute value", http.StatusBadRequest)
		return
	}
	appointmentTime := time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, 0, 0, time.UTC)

	// TODO: validation
	apt := calendar.NewAppointment(0, name, appointmentTime)

	if err = s.state.aptsStorage.Save(apt); err != nil {
		// TODO: check error type
		slog.Error(
			"failed to save appointment",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDeleteAppointment(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("DELETE %s", r.URL))

	aptIDStr := r.PathValue("appointmentID")
	if aptIDStr == "" {
		http.Error(w, "Missing appointmentID path parameter", http.StatusBadRequest)
		return
	}
	aptID, err := strconv.Atoi(aptIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	found, err := s.state.aptsStorage.Delete(r.Context(), aptID)
	if err != nil {
		slog.Error(
			"failed to delete appointment",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, fmt.Sprintf("Appointment %d not found", aptID), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
