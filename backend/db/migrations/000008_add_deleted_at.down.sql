ALTER TABLE rooms DROP COLUMN deleted_at;
ALTER TABLE teams DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;
ALTER TABLE seats DROP COLUMN deleted_at;
ALTER TABLE reservations DROP COLUMN deleted_at;
ALTER TABLE cross_team_requests DROP COLUMN deleted_at;
