export const queryKeys = {
  rooms: {
    all: ['rooms'] as const,
    detail: (id: string) => ['rooms', id] as const,
    seats: (roomId: string) => ['rooms', roomId, 'seats'] as const,
  },
  seats: {
    all: ['seats'] as const,
    detail: (id: string) => ['seats', id] as const,
  },
  teams: {
    all: ['teams'] as const,
    detail: (id: string) => ['teams', id] as const,
  },
  users: {
    all: ['users'] as const,
    detail: (id: string) => ['users', id] as const,
    byEmail: (email: string) => ['users', 'email', email] as const,
  },
  reservations: {
    byDate: (date: string) => ['reservations', date] as const,
    byUserDate: (userId: string, date: string) => ['reservations', userId, date] as const,
  },
  crossTeamRequests: {
    all: ['cross-team-requests'] as const,
    pendingByTeam: (teamId: string) => ['cross-team-requests', 'pending', teamId] as const,
    detail: (id: string) => ['cross-team-requests', id] as const,
  },
} as const;
