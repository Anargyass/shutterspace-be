-- seed/003_seed_users.sql
-- Password untuk semua akun: "Password123!"
-- Hash bcrypt cost 12 dari "Password123!" — generate ulang di produksi!
-- Untuk membuat hash baru: golang.org/x/crypto/bcrypt atau: htpasswd -bnBC 12 "" "Password123!" | tr -d ':\n' | sed 's/$2y/$2a/'

INSERT INTO users (id, name, email, password_hash, phone, role, is_active) VALUES
(
    gen_random_uuid(),
    'Ahmad Nayottama',
    'user@shutterspace.id',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4oQaGKuXFa',
    '+62812345678901',
    'user',
    true
),(
    gen_random_uuid(),
    'Studio Admin',
    'admin@shutterspace.id',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4oQaGKuXFa',
    '+62898765432100',
    'studio_admin',
    true
)
ON CONFLICT (email) DO NOTHING;

-- Assign admin ke semua studio (setelah user admin dibuat)
UPDATE studios
SET managed_by = (SELECT id FROM users WHERE email = 'admin@shutterspace.id')
WHERE managed_by IS NULL;
