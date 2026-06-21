-- 007_create_payments.sql
CREATE TABLE IF NOT EXISTS payments (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID            NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE RESTRICT,
    amount          NUMERIC(12,2)   NOT NULL CHECK (amount >= 0),
    payment_method  payment_method  NOT NULL DEFAULT 'mock_payment',
    status          payment_status  NOT NULL DEFAULT 'pending',
    external_ref    VARCHAR(100),
    paid_at         TIMESTAMPTZ,
    expired_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_booking ON payments (booking_id);
CREATE INDEX IF NOT EXISTS idx_payments_status  ON payments (status);
