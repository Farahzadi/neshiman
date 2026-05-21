package domain

type Position struct {
	X int
	Y int
}

type Rotation int

const (
	Rotation0   Rotation = 0
	Rotation90  Rotation = 90
	Rotation180 Rotation = 180
	Rotation270 Rotation = 270
)

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
