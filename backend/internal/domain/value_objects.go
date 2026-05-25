package domain

type Position struct {
	X int
	Y int
}

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleTeamAdmin  Role = "team_admin"
	RoleViewer     Role = "viewer"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

type WeeklyLimit int

const MaxWeeklyLimit WeeklyLimit = 5
