-- 001_create_enums.sql
-- Jalankan pertama kali sebelum tabel apapun dibuat

CREATE TYPE user_role AS ENUM ('user', 'studio_admin');

CREATE TYPE surabaya_area AS ENUM (
    'Surabaya Pusat',
    'Surabaya Barat',
    'Surabaya Timur',
    'Surabaya Selatan',
    'Surabaya Utara'
);

CREATE TYPE day_of_week AS ENUM (
    'monday', 'tuesday', 'wednesday', 'thursday',
    'friday', 'saturday', 'sunday'
);

CREATE TYPE booking_status AS ENUM (
    'pending', 'confirmed', 'cancelled', 'completed'
);

CREATE TYPE payment_status AS ENUM (
    'pending', 'paid', 'failed', 'refunded'
);

CREATE TYPE payment_method AS ENUM (
    'bank_transfer', 'credit_debit_card', 'e_wallet', 'mock_payment'
);
