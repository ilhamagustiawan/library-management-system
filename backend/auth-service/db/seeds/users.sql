-- Development-only users. Production provisioning intentionally remains external.
-- Both accounts use the password "password".
INSERT INTO users (
    id,
    name,
    email,
    password_hash,
    role_code,
    created_at,
    updated_at
)
VALUES
    (
        '00000000-0000-4000-8000-000000000001',
        'Library Admin',
        'admin@library.com',
        '$2a$12$EHF9iwYY6XPzEjIo59OeoeZxECKVKNHUcSC0XGrJuSMjPD6Mh4dDK',
        'admin',
        CURRENT_TIMESTAMP(6),
        CURRENT_TIMESTAMP(6)
    ),
    (
        '00000000-0000-4000-8000-000000000002',
        'Library Member',
        'member@library.com',
        '$2a$12$gzV3SCnYrfrFrMSnThWXU.bWe3QURoxXuhDj88yQYcHP0rQsxQmKO',
        'member',
        CURRENT_TIMESTAMP(6),
        CURRENT_TIMESTAMP(6)
    ) AS incoming
ON DUPLICATE KEY UPDATE
    name = incoming.name,
    email = incoming.email,
    password_hash = incoming.password_hash,
    role_code = incoming.role_code,
    updated_at = incoming.updated_at;
