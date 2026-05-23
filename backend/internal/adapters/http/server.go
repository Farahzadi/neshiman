package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"neshiman/backend/internal/application"
)

type Server struct {
	server *http.Server
}

func NewServer(
	addr string,
	jwtSecret string,
	corsOrigins string,
	roomSvc *application.RoomService,
	reservationSvc *application.ReservationService,
	teamSvc *application.TeamService,
	seatSvc *application.SeatService,
	userSvc *application.UserService,
	crossTeamRequestSvc *application.CrossTeamRequestService,
	authSvc *application.AuthService,
) *Server {
	r := newRouter(roomSvc, reservationSvc, teamSvc, seatSvc, userSvc, crossTeamRequestSvc, authSvc, jwtSecret, corsOrigins)

	return &Server{
		server: &http.Server{
			Addr:         ":" + addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
