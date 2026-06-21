-- 004_create_studios.sql
CREATE TABLE IF NOT EXISTS studios (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    studio_type_id  INTEGER        NOT NULL REFERENCES studio_types(id) ON DELETE RESTRICT,
    managed_by      UUID           REFERENCES users(id) ON DELETE SET NULL,
    name            VARCHAR(150)   NOT NULL,
    slug            VARCHAR(180)   NOT NULL UNIQUE,
    description     TEXT,
    address         TEXT           NOT NULL,
    area            surabaya_area  NOT NULL,
    capacity        INTEGER        NOT NULL DEFAULT 1 CHECK (capacity > 0),
    area_sqm        NUMERIC(6,2),
    price_per_hour  NUMERIC(12,2)  NOT NULL CHECK (price_per_hour > 0),
    facilities      JSONB          NOT NULL DEFAULT '[]',
    addons          JSONB          NOT NULL DEFAULT '[]',
    images          JSONB          NOT NULL DEFAULT '[]',
    rating          NUMERIC(3,2)   DEFAULT 0.00 CHECK (rating BETWEEN 0 AND 5),
    review_count    INTEGER        NOT NULL DEFAULT 0,
    is_active       BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_studios_type_area  ON studios (studio_type_id, area);
CREATE INDEX IF NOT EXISTS idx_studios_active      ON studios (is_active);
CREATE INDEX IF NOT EXISTS idx_studios_managed_by  ON studios (managed_by);
