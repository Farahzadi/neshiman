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

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	m, err := migrate.New("file://db/migrations", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Printf("warning: migrator source close: %v", srcErr)
	}
	if dbErr != nil {
		log.Printf("warning: migrator database close: %v", dbErr)
	}

	roomRepo := postgres.NewRoomRepository(pool)
	seatRepo := postgres.NewSeatRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	reservationRepo := postgres.NewReservationRepository(pool)
	teamRepo := postgres.NewTeamRepository(pool)
	txManager := postgres.NewTxManager(pool)

	roomSvc := application.NewRoomService(roomRepo, seatRepo)
	seatSvc := application.NewSeatService(seatRepo)
	userSvc := application.NewUserService(userRepo)
	reservationSvc := application.NewReservationService(reservationRepo, seatRepo, userRepo, txManager)
	teamSvc := application.NewTeamService(teamRepo)
	crossTeamRequestRepo := postgres.NewCrossTeamRequestRepository(pool)
	crossTeamRequestSvc := application.NewCrossTeamRequestService(crossTeamRequestRepo, userRepo, seatRepo, reservationRepo, txManager)
	authSvc := application.NewAuthService(userRepo, teamRepo, cfg.JWTSecret)

	srv := httpadapter.NewServer(cfg.Port, cfg.JWTSecret, cfg.CORSOrigins, roomSvc, reservationSvc, teamSvc, seatSvc, userSvc, crossTeamRequestSvc, authSvc)
	fmt.Printf("server listening on :%s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
