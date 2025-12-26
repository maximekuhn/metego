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

// GET /namedays
func (s *Server) namedaysHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /namedays")

	namedays, err := s.state.namedaysStorage.GetAll(r.Context(), 1000, 0)
	if err != nil {
		slog.Error(
			"failed to get namedays from db",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	namedaysPage := views.Namedays(namedays)
	w.Header().Add("Content-Type", "text/html")
	err = namedaysPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render views.Namedays", slog.String("err_msg", err.Error()))
	}
}

// POST /api/namedays
func (s *Server) handleCreateNameday(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST /api/namedays")

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

	nday := calendar.NewNameday(0, name, calendar.BirthdayDate{
		Month: time.Month(m),
		Day:   uint8(d),
	})

	err = s.state.namedaysStorage.Save(r.Context(), nday)
	if err != nil {
		slog.Error("failed to save nday", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDeleteNameday(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("DELETE %s", r.URL))

	ndayIDStr := r.PathValue("namedayID")
	if ndayIDStr == "" {
		http.Error(w, "Missing namedayID path parameter", http.StatusBadRequest)
		return
	}
	ndayID, err := strconv.Atoi(ndayIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	found, err := s.state.namedaysStorage.Delete(r.Context(), ndayID)
	if err != nil {
		slog.Error(
			"failed to delete nameday",
			slog.String("err_msg", err.Error()),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, fmt.Sprintf("Nameday %d not found", ndayID), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
