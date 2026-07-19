-- Development-only clients. Production provisioning intentionally remains external.
INSERT INTO oauth_clients (
    id,
    name,
    kind,
    secret_hash,
    redirect_uri,
    is_public,
    user_id
)
VALUES
    (
        'member-nextjs-web',
        'Library Management System for Members',
        'authorization_code',
        '$2a$12$xu.u61sDGKugXdfRl0e50.ezdFirPQ9eKYOh6e.oXLP2DxItYIeli',
        'http://localhost:3000/api/auth/callback/library',
        FALSE,
        NULL
    ),
    (
        'kong-gateway',
        'Kong Gateway',
        'resource_server',
        '$2a$12$ZRjbNk9m/aJgO9Y9nXqIVufmZ2VrlodAZoxKC.86nb4U9uruQYkZ2',
        '',
        FALSE,
        NULL
    ),
    (
        'book-service',
        'Book Service',
        'resource_server',
        '$2y$12$C.EDWokJc3lDB7d.tIPzO.hMk4Y/FTfZ7ArKUnjx.6xxR//gK3Tma',
        '',
        FALSE,
        NULL
    ),
    (
        'user-service',
        'User Service',
        'client_credentials',
        '$2y$12$f5Njq/MigcimvDxlU.f7EOI5UD4oHRx8PhUZ2.oIpL.bOMZxm.hKK',
        '',
        FALSE,
        NULL
    ),
    (
        'transaction-service',
        'Transaction Service',
        'client_credentials',
        '$2y$12$H97Z6YTm7t1pDQLqitko4.x9wFJeuSPPj7KmBPIJRhMZoQ/ZIKCr.',
        '',
        FALSE,
        NULL
    ) AS incoming
ON DUPLICATE KEY UPDATE
    name = incoming.name,
    kind = incoming.kind,
    secret_hash = incoming.secret_hash,
    redirect_uri = incoming.redirect_uri,
    is_public = incoming.is_public,
    user_id = incoming.user_id;

DELETE FROM oauth_client_scopes
WHERE client_id IN ('member-nextjs-web', 'user-service', 'transaction-service');

-- Human client ceiling includes both roles; role_scopes still prevents privilege escalation.
INSERT INTO oauth_client_scopes (client_id, scope_code) VALUES
    ('member-nextjs-web', 'books:read'),
    ('member-nextjs-web', 'loans:borrow:self'),
    ('member-nextjs-web', 'loans:return:self'),
    ('member-nextjs-web', 'transactions:read:self'),
    ('member-nextjs-web', 'transactions:read:any'),
    ('member-nextjs-web', 'loans:return:any'),
    ('member-nextjs-web', 'fines:manage'),
    ('member-nextjs-web', 'books:manage'),
    ('user-service', 'identities:create'),
    ('transaction-service', 'book-stock:read'),
    ('transaction-service', 'book-stock:reserve'),
    ('transaction-service', 'book-stock:release');
