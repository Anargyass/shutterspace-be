-- seed/002_seed_studios_and_slots.sql

-- Insert 4 studio contoh
INSERT INTO studios (
    id, studio_type_id, name, slug, description, address, area,
    capacity, price_per_hour, facilities, addons, images, rating, review_count
) VALUES
(
    'a1b2c3d4-0001-4001-a001-000000000001',
    (SELECT id FROM studio_types WHERE slug = 'fotografi'),
    'Lumina Studio', 'lumina-studio',
    'Studio foto minimalis dengan pencahayaan natural berlimpah. Backdrop beragam untuk berbagai kebutuhan pemotretan profesional.',
    'Jl. Darmo Permai Selatan No. 12', 'Surabaya Pusat', 8, 150000,
    '["2x Godox Strobes 400W","Octa Softbox","Rectangular Softbox","Backdrop 12 Warna","Make-up Area","WiFi 100Mbps","AC Central","Cermin Besar"]',
    '[{"name":"Seamless Paper Roll","price":50000},{"name":"C-Stand Extra","price":25000},{"name":"Ring Light Additional","price":35000}]',
    '["https://images.unsplash.com/photo-1542038784456-1ea8e935640e?w=800&q=80"]',
    4.8, 124
),(
    'a1b2c3d4-0002-4002-a002-000000000002',
    (SELECT id FROM studio_types WHERE slug = 'podcast'),
    'Echo Pod', 'echo-pod',
    'Studio podcast kedap suara modern. Ideal untuk rekaman solo, interview duo, hingga panel diskusi 4 orang dengan kualitas broadcast.',
    'Jl. Raya Darmo Park III No. 5', 'Surabaya Barat', 4, 120000,
    '["Shure SM7B x4","Focusrite Scarlett 18i20","Acoustic Treatment","Headphone Monitor x4","Layar Monitor Tamu","WiFi 200Mbps","AC","Live Streaming Ready"]',
    '[{"name":"Kopi\/Teh Set","price":25000},{"name":"Tambahan Mic Stand","price":15000},{"name":"Green Screen","price":75000}]',
    '["https://images.unsplash.com/photo-1590602847861-f357a9332bbc?w=800&q=80"]',
    4.7, 89
),(
    'a1b2c3d4-0003-4003-a003-000000000003',
    (SELECT id FROM studio_types WHERE slug = 'produksi-musik'),
    'Resonance Room', 'resonance-room',
    'Studio recording full-equipped untuk band, musisi solo, hingga komposer. Mixing room terpisah dengan acoustic treatment profesional.',
    'Jl. Nginden Semolo No. 45', 'Surabaya Timur', 6, 250000,
    '["Vocal Booth Kedap Suara","Drum Kit Pearl","Bass & Guitar Amp","Piano Digital Yamaha","Pro Tools \/ Logic Pro X","Mixing Console","Monitor Speaker"]',
    '[{"name":"Session Musician","price":200000},{"name":"Recording Engineer","price":150000},{"name":"Mix & Master Session","price":300000}]',
    '["https://images.unsplash.com/photo-1598488035139-bdbb2231ce04?w=800&q=80"]',
    4.9, 56
),(
    'a1b2c3d4-0004-4004-a004-000000000004',
    (SELECT id FROM studio_types WHERE slug = 'fotografi'),
    'Studio Pujas', 'studio-pujas',
    'Studio foto warm & earthy dengan jendela panoramik besar menghadirkan cahaya natural yang melimpah. Ideal untuk editorial, fashion, dan portrait.',
    'Jl. Genteng Kali No. 8', 'Surabaya Pusat', 15, 150000,
    '["2x Godox Strobes 400W","Octa & Rectangular Softboxes","WiFi 100Mbps","Make-up & Fitting Area","AC Central","Bluetooth Speaker Marshall","Ruang Tunggu"]',
    '[{"name":"Seamless Paper Roll","price":50000},{"name":"C-Stand Extra","price":25000},{"name":"Props Dekorasi Set","price":100000}]',
    '["https://images.unsplash.com/photo-1471341971476-ae15ff5dd4ea?w=800&q=80"]',
    4.9, 124
)
ON CONFLICT (slug) DO NOTHING;

-- Insert jam operasional untuk semua studio:
-- Senin-Sabtu: 08:00-22:00 | Minggu: 09:00-18:00
INSERT INTO availability_slots (studio_id, day_of_week, open_time, close_time, is_open)
SELECT
    s.id,
    d.day::day_of_week,
    CASE WHEN d.day = 'sunday' THEN '09:00:00'::TIME ELSE '08:00:00'::TIME END AS open_time,
    CASE WHEN d.day = 'sunday' THEN '18:00:00'::TIME ELSE '22:00:00'::TIME END AS close_time,
    TRUE AS is_open
FROM studios s
CROSS JOIN (VALUES
    ('monday'), ('tuesday'), ('wednesday'), ('thursday'),
    ('friday'), ('saturday'), ('sunday')
) AS d(day)
ON CONFLICT (studio_id, day_of_week) DO NOTHING;
