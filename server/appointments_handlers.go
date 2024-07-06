package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/maximekuhn/metego/calendar"
	"github.com/maximekuhn/metego/server/views"
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
	apt := calendar.NewAppointment(name, appointmentTime)

	if err = s.state.aptsStorage.Save(apt); err != nil {
		// TODO: check error type
		slog.Error(
			"failed to save appointment",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
