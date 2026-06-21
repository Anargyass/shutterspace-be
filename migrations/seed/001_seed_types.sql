-- seed/001_seed_types.sql
INSERT INTO studio_types (name, slug, description) VALUES
    ('Fotografi',      'fotografi',       'Studio pencahayaan profesional untuk foto komersial, portrait, dan fashion.'),
    ('Podcast',        'podcast',         'Studio kedap suara dengan peralatan rekaman audio berkualitas tinggi.'),
    ('Tari',           'tari',            'Studio luas dengan lantai kayu dan cermin penuh untuk latihan koreografi.'),
    ('Produksi Musik', 'produksi-musik',  'Studio recording full-equipped dengan mixing board dan vocal booth.')
ON CONFLICT (slug) DO NOTHING;
