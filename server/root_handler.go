package server

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/metego/server/views"
)

// GET /?city=...
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	city := queryParams.Get("city")
	slog.Info("GET /?city={city}", slog.String("city", city))

	w.Header().Add("Content-Type", "text/html")

	indexPage := views.Index(city)
	err := indexPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Index page", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GET /admin
func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /admin")

	adminPage := views.Admin()
	err := adminPage.Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Index page", slog.String("err_msg", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
