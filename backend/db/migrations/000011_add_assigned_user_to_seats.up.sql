ALTER TABLE seats ADD COLUMN assigned_user_id UUID REFERENCES users(id);

CREATE UNIQUE INDEX idx_seats_assigned_user ON seats(assigned_user_id)
  WHERE assigned_user_id IS NOT NULL AND deleted_at IS NULL;
