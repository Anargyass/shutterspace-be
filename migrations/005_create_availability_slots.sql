-- 005_create_availability_slots.sql
-- Menyimpan jam operasional studio per hari dalam seminggu
CREATE TABLE IF NOT EXISTS availability_slots (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    studio_id   UUID         NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    day_of_week day_of_week  NOT NULL,
    open_time   TIME         NOT NULL DEFAULT '08:00:00',
    close_time  TIME         NOT NULL DEFAULT '22:00:00',
    is_open     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_studio_day  UNIQUE (studio_id, day_of_week),
    CONSTRAINT chk_op_hours   CHECK  (close_time > open_time)
);

CREATE INDEX IF NOT EXISTS idx_slots_studio ON availability_slots (studio_id);
