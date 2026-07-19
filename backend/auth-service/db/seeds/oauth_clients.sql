-- Development-only clients. Production provisioning intentionally remains external.
INSERT INTO oauth_clients (
    id,
    name,
    kind,
    secret_hash,
    redirect_uri,
    allowed_scopes,
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
        'library:read library:write',
        FALSE,
        NULL
    ),
    (
        'kong-gateway',
        'Kong Gateway',
        'resource_server',
        '$2a$12$ZRjbNk9m/aJgO9Y9nXqIVufmZ2VrlodAZoxKC.86nb4U9uruQYkZ2',
        '',
        '',
        FALSE,
        NULL
    ) AS incoming
ON DUPLICATE KEY UPDATE
    name = incoming.name,
    kind = incoming.kind,
    secret_hash = incoming.secret_hash,
    redirect_uri = incoming.redirect_uri,
    allowed_scopes = incoming.allowed_scopes,
    is_public = incoming.is_public,
    user_id = incoming.user_id;
