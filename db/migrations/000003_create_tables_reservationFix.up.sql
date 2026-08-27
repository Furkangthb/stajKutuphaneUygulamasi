ALTER TABLE reservations DROP CONSTRAINT reservations_status_check;
ALTER TABLE reservations ADD CONSTRAINT reservations_status_check
    CHECK (status IN ('pending', 'active', 'completed', 'cancelled', 'expired'));
ALTER TABLE reservations ALTER COLUMN status SET DEFAULT 'pending';