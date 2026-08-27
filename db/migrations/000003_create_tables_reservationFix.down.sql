ALTER TABLE reservations ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE reservations DROP CONSTRAINT reservations_status_check;
ALTER TABLE reservations ADD CONSTRAINT reservations_status_check
    CHECK (status IN ('active', 'completed', 'cancelled', 'expired'));