package application

import (
	"context"
	"neshiman/backend/internal/domain"
	"neshiman/backend/internal/ports"

	"github.com/google/uuid"
)

type mockRoomRepo struct {
	createFn func(ctx context.Context, room *domain.Room) error
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Room, error)
	listFn   func(ctx context.Context) ([]domain.Room, error)
	updateFn func(ctx context.Context, room *domain.Room) error
	deleteFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockRoomRepo) Create(ctx context.Context, room *domain.Room) error { return m.createFn(ctx, room) }
func (m *mockRoomRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) { return m.getByIDFn(ctx, id) }
func (m *mockRoomRepo) List(ctx context.Context) ([]domain.Room, error) { return m.listFn(ctx) }
func (m *mockRoomRepo) Update(ctx context.Context, room *domain.Room) error { return m.updateFn(ctx, room) }
func (m *mockRoomRepo) Delete(ctx context.Context, id uuid.UUID) error { return m.deleteFn(ctx, id) }
func (m *mockRoomRepo) assertImplementation() { var _ ports.RoomRepository = m }

type mockSeatRepo struct {
	createFn      func(ctx context.Context, seat *domain.Seat) error
	getByIDFn     func(ctx context.Context, id uuid.UUID) (*domain.Seat, error)
	listByRoomFn  func(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error)
	updateFn      func(ctx context.Context, seat *domain.Seat) error
	deleteFn      func(ctx context.Context, id uuid.UUID) error
	deleteByRoomFn func(ctx context.Context, roomID uuid.UUID) error
	bulkSyncFn    func(ctx context.Context, roomID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error)
}

func (m *mockSeatRepo) Create(ctx context.Context, seat *domain.Seat) error { return m.createFn(ctx, seat) }
func (m *mockSeatRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Seat, error) { return m.getByIDFn(ctx, id) }
func (m *mockSeatRepo) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Seat, error) { return m.listByRoomFn(ctx, roomID) }
func (m *mockSeatRepo) Update(ctx context.Context, seat *domain.Seat) error { return m.updateFn(ctx, seat) }
func (m *mockSeatRepo) Delete(ctx context.Context, id uuid.UUID) error { return m.deleteFn(ctx, id) }
func (m *mockSeatRepo) DeleteByRoom(ctx context.Context, roomID uuid.UUID) error { return m.deleteByRoomFn(ctx, roomID) }
func (m *mockSeatRepo) BulkSync(ctx context.Context, roomID uuid.UUID, seats []domain.Seat) ([]domain.Seat, error) { return m.bulkSyncFn(ctx, roomID, seats) }
func (m *mockSeatRepo) assertImplementation() { var _ ports.SeatRepository = m }

type mockUserRepo struct {
	createFn            func(ctx context.Context, user *domain.User) error
	getByIDFn           func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	getByEmailFn        func(ctx context.Context, email string) (*domain.User, error)
	getByNameFn         func(ctx context.Context, name string) (*domain.User, error)
	listByTeamFn        func(ctx context.Context, teamID uuid.UUID) ([]domain.User, error)
	listAllFn           func(ctx context.Context) ([]domain.User, error)
	updateWeeklyLimitFn func(ctx context.Context, id uuid.UUID, limit domain.WeeklyLimit) error
	updatePasswordFn    func(ctx context.Context, id uuid.UUID, passwordHash string) error
	deleteFn            func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error { return m.createFn(ctx, user) }
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) { return m.getByIDFn(ctx, id) }
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) { return m.getByEmailFn(ctx, email) }
func (m *mockUserRepo) GetByName(ctx context.Context, name string) (*domain.User, error) { return m.getByNameFn(ctx, name) }
func (m *mockUserRepo) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.User, error) { return m.listByTeamFn(ctx, teamID) }
func (m *mockUserRepo) ListAll(ctx context.Context) ([]domain.User, error) { return m.listAllFn(ctx) }
func (m *mockUserRepo) UpdateWeeklyLimit(ctx context.Context, id uuid.UUID, limit domain.WeeklyLimit) error { return m.updateWeeklyLimitFn(ctx, id, limit) }
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error { return m.updatePasswordFn(ctx, id, passwordHash) }
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return m.deleteFn(ctx, id) }

type mockReservationRepo struct {
	createFn           func(ctx context.Context, r *domain.Reservation) error
	getByIDFn          func(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	getBySeatAndDateFn   func(ctx context.Context, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error)
	listByDateFn       func(ctx context.Context, date domain.Date) ([]domain.Reservation, error)
	listByUserAndDateFn func(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error)
	listByUserAndWeekFn func(ctx context.Context, userID uuid.UUID, start, end domain.Date) ([]domain.Reservation, error)
	countByUserInWeekFn func(ctx context.Context, userID uuid.UUID, start, end domain.Date) (int, error)
	deleteFn           func(ctx context.Context, id uuid.UUID) error
}

func (m *mockReservationRepo) Create(ctx context.Context, r *domain.Reservation) error { return m.createFn(ctx, r) }
func (m *mockReservationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) { return m.getByIDFn(ctx, id) }
func (m *mockReservationRepo) GetBySeatAndDate(ctx context.Context, seatID uuid.UUID, date domain.Date) (*domain.Reservation, error) { return m.getBySeatAndDateFn(ctx, seatID, date) }
func (m *mockReservationRepo) ListByDate(ctx context.Context, date domain.Date) ([]domain.Reservation, error) { return m.listByDateFn(ctx, date) }
func (m *mockReservationRepo) ListByUserAndDate(ctx context.Context, userID uuid.UUID, date domain.Date) ([]domain.Reservation, error) { return m.listByUserAndDateFn(ctx, userID, date) }
func (m *mockReservationRepo) ListByUserAndWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) ([]domain.Reservation, error) { return m.listByUserAndWeekFn(ctx, userID, start, end) }
func (m *mockReservationRepo) CountByUserInWeek(ctx context.Context, userID uuid.UUID, start, end domain.Date) (int, error) { return m.countByUserInWeekFn(ctx, userID, start, end) }
func (m *mockReservationRepo) Delete(ctx context.Context, id uuid.UUID) error { return m.deleteFn(ctx, id) }

type mockTeamRepo struct {
	createFn func(ctx context.Context, team *domain.Team) error
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	listFn   func(ctx context.Context) ([]domain.Team, error)
	deleteFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTeamRepo) Create(ctx context.Context, team *domain.Team) error { return m.createFn(ctx, team) }
func (m *mockTeamRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) { return m.getByIDFn(ctx, id) }
func (m *mockTeamRepo) List(ctx context.Context) ([]domain.Team, error) { return m.listFn(ctx) }
func (m *mockTeamRepo) Delete(ctx context.Context, id uuid.UUID) error { return m.deleteFn(ctx, id) }

type mockCrossTeamRequestRepo struct {
	createFn         func(ctx context.Context, r *domain.CrossTeamRequest) error
	getByIDFn        func(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error)
	listByStatusFn   func(ctx context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error)
	listPendingByTeamFn func(ctx context.Context, teamID uuid.UUID) ([]domain.CrossTeamRequest, error)
	updateStatusFn   func(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error
}

func (m *mockCrossTeamRequestRepo) Create(ctx context.Context, r *domain.CrossTeamRequest) error { return m.createFn(ctx, r) }
func (m *mockCrossTeamRequestRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CrossTeamRequest, error) { return m.getByIDFn(ctx, id) }
func (m *mockCrossTeamRequestRepo) ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.CrossTeamRequest, error) { return m.listByStatusFn(ctx, status) }
func (m *mockCrossTeamRequestRepo) ListPendingByTeam(ctx context.Context, teamID uuid.UUID) ([]domain.CrossTeamRequest, error) { return m.listPendingByTeamFn(ctx, teamID) }
func (m *mockCrossTeamRequestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error { return m.updateStatusFn(ctx, id, status) }

type mockTxManager struct {
	withinTxFn func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *mockTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.withinTxFn(ctx, fn)
}
