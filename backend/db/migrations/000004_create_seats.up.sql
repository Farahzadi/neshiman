CREATE TABLE seats (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id   UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    team_id   UUID NOT NULL REFERENCES teams(id),
    label     TEXT NOT NULL,
    pos_x     INT NOT NULL,
    pos_y     INT NOT NULL,
    rotation  INT NOT NULL DEFAULT 0 CHECK (rotation IN (0, 90, 180, 270)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (room_id, pos_x, pos_y)
);
