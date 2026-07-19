CREATE TABLE users (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    password_hash VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    UNIQUE KEY users_email_uidx (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE sessions (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    user_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    token_hash CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    expires_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    UNIQUE KEY sessions_token_hash_uidx (token_hash),
    KEY sessions_user_id_idx (user_id),
    KEY sessions_expires_at_idx (expires_at),
    CONSTRAINT sessions_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE oauth_clients (
    id VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    secret_hash VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    redirect_uri VARCHAR(2048) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    allowed_scopes TEXT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    user_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    KEY oauth_clients_user_id_idx (user_id),
    CONSTRAINT oauth_clients_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT oauth_clients_redirect_uri_nonempty CHECK (redirect_uri <> '')
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE oauth_tokens (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    code VARCHAR(255) CHARACTER SET ascii COLLATE ascii_bin NULL,
    access_token VARCHAR(255) CHARACTER SET ascii COLLATE ascii_bin NULL,
    refresh_token VARCHAR(255) CHARACTER SET ascii COLLATE ascii_bin NULL,
    payload JSON NOT NULL,
    expires_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    UNIQUE KEY oauth_tokens_code_uidx (code),
    UNIQUE KEY oauth_tokens_access_uidx (access_token),
    UNIQUE KEY oauth_tokens_refresh_uidx (refresh_token),
    KEY oauth_tokens_expires_at_idx (expires_at),
    CONSTRAINT oauth_tokens_code_or_access CHECK (
        (code IS NOT NULL AND access_token IS NULL)
        OR (code IS NULL AND access_token IS NOT NULL)
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
