package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"neshiman/backend/internal/adapters/config"
	"neshiman/backend/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	q := sqlc.New(pool)

	fmt.Println("Seeding database...")
	truncateAll(ctx, pool)

	teams := seedTeams(ctx, q)
	users := seedUsers(ctx, q, teams)
	rooms := seedRooms(ctx, q)
	seats := seedSeats(ctx, q, rooms, teams)
	seedReservations(ctx, q, users, seats)
	seedCrossTeamRequests(ctx, q, users, seats)

	fmt.Println("Done.")
}

func seedTeams(ctx context.Context, q *sqlc.Queries) map[string]sqlc.Team {
	names := []struct {
		key  string
		name string
	}{
		{"eng", "Engineering"},
		{"design", "Design"},
		{"marketing", "Marketing"},
	}

	result := make(map[string]sqlc.Team, len(names))
	for _, n := range names {
		t, err := q.CreateTeam(ctx, n.name)
		if err != nil {
			log.Fatalf("failed to create team %s: %v", n.name, err)
		}
		result[n.key] = t
		fmt.Printf("  Team: %s (id=%s)\n", t.Name, t.ID)
	}
	return result
}

func seedUsers(ctx context.Context, q *sqlc.Queries, teams map[string]sqlc.Team) map[string]sqlc.User {
	type userSeed struct {
		key         string
		name        string
		email       string
		teamKey     string
		role        string
		weeklyLimit int32
	}

	seeds := []userSeed{
		{key: "super", name: "Admin", email: "admin@neshiman.test", teamKey: "", role: "superadmin", weeklyLimit: 5},
		{key: "eng_lead", name: "Alice", email: "alice@neshiman.test", teamKey: "eng", role: "team_admin", weeklyLimit: 3},
		{key: "eng_dev1", name: "Bob", email: "bob@neshiman.test", teamKey: "eng", role: "viewer", weeklyLimit: 2},
		{key: "eng_dev2", name: "Charlie", email: "charlie@neshiman.test", teamKey: "eng", role: "viewer", weeklyLimit: 2},
		{key: "design_lead", name: "Diana", email: "diana@neshiman.test", teamKey: "design", role: "team_admin", weeklyLimit: 3},
		{key: "design_dev1", name: "Eve", email: "eve@neshiman.test", teamKey: "design", role: "viewer", weeklyLimit: 2},
		{key: "mktg_lead", name: "Frank", email: "frank@neshiman.test", teamKey: "marketing", role: "team_admin", weeklyLimit: 3},
		{key: "mktg_dev1", name: "Grace", email: "grace@neshiman.test", teamKey: "marketing", role: "viewer", weeklyLimit: 2},
	}

	result := make(map[string]sqlc.User, len(seeds))
	for _, s := range seeds {
		var teamID pgtype.UUID
		if s.teamKey != "" {
			teamID = pgtype.UUID{Bytes: teams[s.teamKey].ID, Valid: true}
		}
		u, err := q.CreateUser(ctx, sqlc.CreateUserParams{
			Name:  s.name,
			Email: pgtype.Text{String: s.email, Valid: s.email != ""},
			TeamID:      teamID,
			Role:        s.role,
			WeeklyLimit: s.weeklyLimit,
		})
		if err != nil {
			log.Fatalf("failed to create user %s: %v", s.name, err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password for %s: %v", s.name, err)
		}
		if err := q.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
			ID:           u.ID,
			PasswordHash: string(hash),
		}); err != nil {
			log.Fatalf("failed to set password for %s: %v", s.name, err)
		}

		result[s.key] = u
		fmt.Printf("  User: %s (id=%s, team=%q)\n", u.Name, u.ID, s.teamKey)
	}
	return result
}

func seedRooms(ctx context.Context, q *sqlc.Queries) map[string]sqlc.Room {
	type roomSeed struct {
		key        string
		name       string
		gridWidth  int32
		gridHeight int32
	}

	seeds := []roomSeed{
		{key: "conf_a", name: "Conference Room A", gridWidth: 10, gridHeight: 8},
		{key: "open_space", name: "Open Space", gridWidth: 20, gridHeight: 12},
	}

	result := make(map[string]sqlc.Room, len(seeds))
	for _, s := range seeds {
		r, err := q.CreateRoom(ctx, sqlc.CreateRoomParams{
			Name:       s.name,
			GridWidth:  s.gridWidth,
			GridHeight: s.gridHeight,
		})
		if err != nil {
			log.Fatalf("failed to create room %s: %v", s.name, err)
		}
		result[s.key] = r
		fmt.Printf("  Room: %s (id=%s, %dx%d)\n", r.Name, r.ID, r.GridWidth, r.GridHeight)
	}
	return result
}

func seedSeats(ctx context.Context, q *sqlc.Queries, rooms map[string]sqlc.Room, teams map[string]sqlc.Team) map[string]sqlc.Seat {
	type seatSeed struct {
		key     string
		roomKey string
		teamKey string
		label   string
		posX    int32
		posY    int32
	}

	seeds := []seatSeed{
		{key: "conf_a_1", roomKey: "conf_a", teamKey: "eng", label: "A1", posX: 1, posY: 1},
		{key: "conf_a_2", roomKey: "conf_a", teamKey: "eng", label: "A2", posX: 2, posY: 1},
		{key: "conf_a_3", roomKey: "conf_a", teamKey: "design", label: "B1", posX: 1, posY: 2},
		{key: "conf_a_4", roomKey: "conf_a", teamKey: "design", label: "B2", posX: 2, posY: 2},
		{key: "open_1", roomKey: "open_space", teamKey: "eng", label: "E1", posX: 2, posY: 2},
		{key: "open_2", roomKey: "open_space", teamKey: "eng", label: "E2", posX: 4, posY: 2},
		{key: "open_3", roomKey: "open_space", teamKey: "design", label: "D1", posX: 2, posY: 6},
		{key: "open_4", roomKey: "open_space", teamKey: "design", label: "D2", posX: 4, posY: 6},
		{key: "open_5", roomKey: "open_space", teamKey: "marketing", label: "M1", posX: 2, posY: 10},
		{key: "open_6", roomKey: "open_space", teamKey: "marketing", label: "M2", posX: 4, posY: 10},
	}

	result := make(map[string]sqlc.Seat, len(seeds))
	for _, s := range seeds {
		seat, err := q.CreateSeat(ctx, sqlc.CreateSeatParams{
			RoomID: rooms[s.roomKey].ID,
			TeamID: teams[s.teamKey].ID,
			Label:  s.label,
			PosX:   s.posX,
			PosY:   s.posY,
		})
		if err != nil {
			log.Fatalf("failed to create seat %s: %v", s.label, err)
		}
		result[s.key] = seat
		fmt.Printf("  Seat: %s (id=%s, in %s, team=%s)\n", seat.Label, seat.ID, s.roomKey, s.teamKey)
	}
	return result
}

func seedReservations(ctx context.Context, q *sqlc.Queries, users map[string]sqlc.User, seats map[string]sqlc.Seat) {
	type resSeed struct {
		userKey string
		seatKey string
		date    pgtype.Date
	}

	seeds := []resSeed{
		{userKey: "eng_dev1", seatKey: "conf_a_1", date: today(0)},
		{userKey: "eng_dev1", seatKey: "conf_a_1", date: today(1)},
		{userKey: "eng_dev2", seatKey: "conf_a_2", date: today(0)},
		{userKey: "design_dev1", seatKey: "conf_a_3", date: today(0)},
		{userKey: "mktg_dev1", seatKey: "open_5", date: today(0)},
		{userKey: "mktg_dev1", seatKey: "open_5", date: today(1)},
	}

	for _, s := range seeds {
		r, err := q.CreateReservation(ctx, sqlc.CreateReservationParams{
			UserID: users[s.userKey].ID,
			SeatID: seats[s.seatKey].ID,
			Date:   s.date,
		})
		if err != nil {
			log.Fatalf("failed to create reservation: %v", err)
		}
		fmt.Printf("  Reservation: %s -> %s (%s)\n",
			s.userKey, s.seatKey, fmtDate(r.Date))
	}
}

func seedCrossTeamRequests(ctx context.Context, q *sqlc.Queries, users map[string]sqlc.User, seats map[string]sqlc.Seat) {
	type ctrSeed struct {
		userKey string
		seatKey string
		date    pgtype.Date
		status  string
	}

	seeds := []ctrSeed{
		{userKey: "design_dev1", seatKey: "conf_a_1", date: today(3), status: "pending"},
		{userKey: "mktg_dev1", seatKey: "conf_a_2", date: today(3), status: "approved"},
		{userKey: "eng_dev1", seatKey: "open_5", date: today(4), status: "rejected"},
	}

	for _, s := range seeds {
		r, err := q.CreateCrossTeamRequest(ctx, sqlc.CreateCrossTeamRequestParams{
			RequestingUserID: users[s.userKey].ID,
			TargetSeatID:     seats[s.seatKey].ID,
			Date:             s.date,
		})
		if err != nil {
			log.Fatalf("failed to create cross-team request: %v", err)
		}
		if s.status != "pending" {
			_, err := q.UpdateCrossTeamRequestStatus(ctx, sqlc.UpdateCrossTeamRequestStatusParams{
				ID:     r.ID,
				Status: s.status,
			})
			if err != nil {
				log.Fatalf("failed to update cross-team request status: %v", err)
			}
		}
		fmt.Printf("  CrossTeamRequest: %s -> %s (%s)\n", s.userKey, s.seatKey, s.status)
	}
}

func today(offsetDays int) pgtype.Date {
	t := time.Now().UTC().AddDate(0, 0, offsetDays)
	return pgtype.Date{Time: t, Valid: true}
}

func truncateAll(ctx context.Context, pool *pgxpool.Pool) {
	tables := []string{"cross_team_requests", "reservations", "seats", "users", "rooms", "teams"}
	for _, t := range tables {
		if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t)); err != nil {
			log.Fatalf("failed to truncate %s: %v", t, err)
		}
	}
}

func fmtDate(d pgtype.Date) string {
	if !d.Valid {
		return "<nil>"
	}
	return d.Time.Format("2006-01-02")
}
