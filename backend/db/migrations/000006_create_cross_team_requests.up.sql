CREATE TABLE cross_team_requests (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requesting_user_id UUID NOT NULL REFERENCES users(id),
    target_seat_id    UUID NOT NULL REFERENCES seats(id),
    date             DATE NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (target_seat_id, date, requesting_user_id)
);
