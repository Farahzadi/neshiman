ALTER TABLE reservations DROP CONSTRAINT reservations_seat_id_date_key;
CREATE UNIQUE INDEX idx_reservations_seat_date_active ON reservations(seat_id, date) WHERE deleted_at IS NULL;

ALTER TABLE cross_team_requests DROP CONSTRAINT cross_team_requests_target_seat_id_date_requesting_user_id_key;
CREATE UNIQUE INDEX idx_cross_team_requests_active ON cross_team_requests(target_seat_id, date, requesting_user_id) WHERE deleted_at IS NULL;
