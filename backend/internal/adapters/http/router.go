package http

import (
	"neshiman/backend/internal/adapters/http/handlers"
	"neshiman/backend/internal/adapters/http/middleware"
	"neshiman/backend/internal/application"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
)

func newRouter(
	roomSvc *application.RoomService,
	reservationSvc *application.ReservationService,
	teamSvc *application.TeamService,
	seatSvc *application.SeatService,
	userSvc *application.UserService,
	crossTeamRequestSvc *application.CrossTeamRequestService,
	authSvc *application.AuthService,
	jwtSecret string,
	corsOrigins string,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.CORS(corsOrigins))
	r.Use(middleware.RequestID)

	r.Get("/health", handlers.Health)

	r.Get("/swagger/*", httpSwagger.Handler())

	roomHandler := handlers.NewRoomHandler(roomSvc)
	reservationHandler := handlers.NewReservationHandler(reservationSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	seatHandler := handlers.NewSeatHandler(seatSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	crossTeamRequestHandler := handlers.NewCrossTeamRequestHandler(crossTeamRequestSvc)
	authHandler := handlers.NewAuthHandler(authSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret))

			r.Route("/rooms", func(r chi.Router) {
				r.Post("/", roomHandler.Create)
				r.Get("/", roomHandler.List)
				r.Get("/{id}", roomHandler.GetByID)
				r.Put("/{id}", roomHandler.Update)
				r.Put("/{id}/seats", seatHandler.BulkSync)
				r.Delete("/{id}", roomHandler.Delete)
			})

			r.Route("/seats", func(r chi.Router) {
				r.Post("/", seatHandler.Create)
				r.Get("/", seatHandler.ListByRoom)
				r.Get("/{id}", seatHandler.GetByID)
				r.Put("/{id}/move", seatHandler.Move)
				r.Put("/{id}/rotate", seatHandler.Rotate)
				r.Delete("/{id}", seatHandler.Delete)
			})

			r.Route("/teams", func(r chi.Router) {
				r.Post("/", teamHandler.Create)
				r.Get("/", teamHandler.List)
				r.Get("/{id}", teamHandler.GetByID)
				r.Delete("/{id}", teamHandler.Delete)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.ListByTeam)
				r.Get("/by-email", userHandler.GetByEmail)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}/password", authHandler.SetPassword)
				r.Put("/{id}/weekly-limit", userHandler.UpdateWeeklyLimit)
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole("team_admin"))
					r.Post("/", userHandler.Create)
					r.Delete("/{id}", userHandler.Delete)
				})
			})

			r.Route("/cross-team-requests", func(r chi.Router) {
				r.Post("/", crossTeamRequestHandler.Create)
				r.Get("/", crossTeamRequestHandler.ListByStatus)
				r.Get("/pending-by-team", crossTeamRequestHandler.ListPendingByTeam)
				r.Get("/{id}", crossTeamRequestHandler.GetByID)
				r.Put("/{id}/approve", crossTeamRequestHandler.Approve)
				r.Put("/{id}/reject", crossTeamRequestHandler.Reject)
			})

			r.Route("/reservations", func(r chi.Router) {
				r.Get("/", reservationHandler.List)
				r.Post("/", reservationHandler.Create)
				r.Get("/week", reservationHandler.GetWeek)
				r.Delete("/{id}", reservationHandler.Cancel)
			})
		})
	})

	return r
}
