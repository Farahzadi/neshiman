CREATE TABLE reservations (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    seat_id UUID NOT NULL REFERENCES seats(id),
    date    DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (seat_id, date)
);
