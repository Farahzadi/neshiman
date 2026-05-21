package http

import (
	"neshiman/backend/internal/adapters/http/handlers"
	"neshiman/backend/internal/application"

	"github.com/go-chi/chi/v5"
)

func newRouter(
	roomSvc *application.RoomService,
	reservationSvc *application.ReservationService,
) *chi.Mux {
	r := chi.NewRouter()

	roomHandler := handlers.NewRoomHandler(roomSvc)
	reservationHandler := handlers.NewReservationHandler(reservationSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", roomHandler.Create)
			r.Get("/", roomHandler.List)
			r.Get("/{id}", roomHandler.GetByID)
			r.Put("/{id}", roomHandler.Update)
			r.Delete("/{id}", roomHandler.Delete)
		})

		r.Route("/reservations", func(r chi.Router) {
			r.Post("/", reservationHandler.Create)
			r.Get("/week", reservationHandler.GetWeek)
			r.Delete("/{id}", reservationHandler.Cancel)
		})
	})

	return r
}
