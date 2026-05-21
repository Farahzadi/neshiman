package main

import (
	"context"
	"fmt"
	"log"

	_ "neshiman/backend/docs"
	"neshiman/backend/internal/adapters/config"
	httpadapter "neshiman/backend/internal/adapters/http"
	"neshiman/backend/internal/adapters/postgres"
	"neshiman/backend/internal/application"
)

// @title           Neshiman API
// @version         1.0
// @description     Seat plan management system
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	roomRepo := postgres.NewRoomRepository(pool)
	seatRepo := postgres.NewSeatRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	reservationRepo := postgres.NewReservationRepository(pool)
	teamRepo := postgres.NewTeamRepository(pool)
	txManager := postgres.NewTxManager(pool)

	roomSvc := application.NewRoomService(roomRepo)
	seatSvc := application.NewSeatService(seatRepo)
	userSvc := application.NewUserService(userRepo)
	reservationSvc := application.NewReservationService(reservationRepo, seatRepo, userRepo, txManager)
	teamSvc := application.NewTeamService(teamRepo)
	crossTeamRequestRepo := postgres.NewCrossTeamRequestRepository(pool)
	crossTeamRequestSvc := application.NewCrossTeamRequestService(crossTeamRequestRepo)

	srv := httpadapter.NewServer(cfg.Port, cfg.JWTSecret, roomSvc, reservationSvc, teamSvc, seatSvc, userSvc, crossTeamRequestSvc)
	fmt.Printf("server listening on :%s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
