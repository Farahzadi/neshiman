CREATE TABLE users (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name      TEXT NOT NULL,
    email     TEXT NOT NULL UNIQUE,
    team_id   UUID REFERENCES teams(id),
    role      TEXT NOT NULL CHECK (role IN ('superadmin', 'team_admin', 'viewer')),
    weekly_limit INT NOT NULL DEFAULT 2,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
