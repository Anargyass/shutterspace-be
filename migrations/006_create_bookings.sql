-- 006_create_bookings.sql
CREATE TABLE IF NOT EXISTS bookings (
    id               UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID           NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    studio_id        UUID           NOT NULL REFERENCES studios(id) ON DELETE RESTRICT,
    booking_date     DATE           NOT NULL,
    start_time       TIME           NOT NULL,
    end_time         TIME           NOT NULL,
    duration_hours   NUMERIC(4,2)   NOT NULL CHECK (duration_hours >= 1),
    status           booking_status NOT NULL DEFAULT 'pending',
    price_per_hour   NUMERIC(12,2)  NOT NULL,
    addons_cost      NUMERIC(12,2)  NOT NULL DEFAULT 0,
    service_fee      NUMERIC(12,2)  NOT NULL DEFAULT 0,
    total_amount     NUMERIC(12,2)  NOT NULL CHECK (total_amount >= 0),
    selected_addons  JSONB          NOT NULL DEFAULT '[]',
    notes            TEXT,
    cancelled_at     TIMESTAMPTZ,
    cancel_reason    TEXT,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_booking_time CHECK (end_time > start_time)
);

-- Index untuk query overlap check (paling sering dipanggil)
CREATE INDEX IF NOT EXISTS idx_bookings_studio_date ON bookings (studio_id, booking_date);
CREATE INDEX IF NOT EXISTS idx_bookings_user        ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status      ON bookings (status);

-- Partial unique index: cegah exact start_time conflict untuk booking aktif
CREATE UNIQUE INDEX IF NOT EXISTS uq_no_exact_start_conflict
    ON bookings (studio_id, booking_date, start_time)
    WHERE status IN ('pending', 'confirmed');
