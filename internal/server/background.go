package server

import (
	"context"
	"log/slog"
	"time"
)

func (s *Server) cleanupOldAppointments() {
	for {
		slog.Info("Clean up old appointments...")
		now := time.Now()
		deleted, err := s.state.aptsStorage.DeleteAllBefore(context.Background(), now)
		if err != nil {
			slog.Error("Failed to delete old appointments", slog.String("err_msg", err.Error()))
		} else {
			slog.Info("All old appointments have been cleaned up", slog.Int("deleted", deleted))
		}
		time.Sleep(24 * time.Hour)
	}
}
