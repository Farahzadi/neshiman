DROP INDEX idx_reservations_seat_date_active;
ALTER TABLE reservations ADD CONSTRAINT reservations_seat_id_date_key UNIQUE (seat_id, date);

DROP INDEX idx_cross_team_requests_active;
ALTER TABLE cross_team_requests ADD CONSTRAINT cross_team_requests_target_seat_id_date_requesting_user_id_key UNIQUE (target_seat_id, date, requesting_user_id);
