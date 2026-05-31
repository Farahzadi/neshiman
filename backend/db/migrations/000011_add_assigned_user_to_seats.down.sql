DROP INDEX IF EXISTS idx_seats_assigned_user;
ALTER TABLE seats DROP COLUMN IF EXISTS assigned_user_id;
