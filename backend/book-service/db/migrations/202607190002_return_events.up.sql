CREATE TABLE message_inbox (
    event_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    event_type VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    processed_at DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE message_outbox (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    event_type VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    routing_key VARCHAR(150) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    payload JSON NOT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'pending',
    attempts INT UNSIGNED NOT NULL DEFAULT 0,
    next_attempt_at DATETIME(6) NOT NULL,
    last_error VARCHAR(1000) NULL,
    created_at DATETIME(6) NOT NULL,
    published_at DATETIME(6) NULL,
    KEY message_outbox_pending_idx (status, next_attempt_at, created_at),
    CONSTRAINT message_outbox_status_valid CHECK (status IN ('pending', 'published'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
