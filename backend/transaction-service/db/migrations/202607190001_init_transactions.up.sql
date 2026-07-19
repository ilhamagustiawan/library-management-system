CREATE TABLE member_loan_counters (
    member_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    active_count TINYINT UNSIGNED NOT NULL DEFAULT 0,
    updated_at DATETIME(6) NOT NULL,
    CONSTRAINT member_loan_counters_range CHECK (active_count <= 3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE loans (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    member_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    book_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    status VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    stock_sync_status VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    borrowed_at DATETIME(6) NOT NULL,
    due_at DATETIME(6) NOT NULL,
    returned_at DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    active_marker TINYINT AS (
        CASE WHEN status IN ('pending_reservation', 'active') THEN 1 ELSE NULL END
    ) STORED,
    UNIQUE KEY loans_member_book_active_uidx (member_id, book_id, active_marker),
    KEY loans_member_history_idx (member_id, created_at DESC),
    KEY loans_stock_sync_idx (stock_sync_status, updated_at),
    CONSTRAINT loans_status_valid CHECK (status IN ('pending_reservation', 'active', 'returned', 'cancelled')),
    CONSTRAINT loans_stock_sync_valid CHECK (stock_sync_status IN ('not_applicable', 'pending', 'confirmed')),
    CONSTRAINT loans_dates_valid CHECK (due_at > borrowed_at),
    CONSTRAINT loans_return_state_valid CHECK (
        (status = 'returned' AND returned_at IS NOT NULL AND stock_sync_status IN ('pending', 'confirmed'))
        OR (status <> 'returned' AND returned_at IS NULL AND stock_sync_status = 'not_applicable')
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE fines (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    loan_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    member_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    overdue_days INT UNSIGNED NOT NULL,
    daily_rate_minor BIGINT UNSIGNED NOT NULL,
    total_amount_minor BIGINT UNSIGNED NOT NULL,
    currency CHAR(3) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    assessed_at DATETIME(6) NOT NULL,
    UNIQUE KEY fines_loan_uidx (loan_id),
    KEY fines_member_status_idx (member_id, status),
    CONSTRAINT fines_loan_fk FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fines_values_valid CHECK (
        overdue_days > 0 AND daily_rate_minor > 0
        AND total_amount_minor = overdue_days * daily_rate_minor
        AND currency = 'IDR' AND status = 'unpaid'
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE loan_transactions (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    loan_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    member_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    book_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    type VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    occurred_at DATETIME(6) NOT NULL,
    UNIQUE KEY loan_transactions_loan_type_uidx (loan_id, type),
    KEY loan_transactions_member_time_idx (member_id, occurred_at DESC),
    KEY loan_transactions_time_idx (occurred_at DESC),
    CONSTRAINT loan_transactions_loan_fk FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT loan_transactions_type_valid CHECK (type IN ('borrow', 'return'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE outbox_events (
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
    KEY outbox_pending_idx (status, next_attempt_at, created_at),
    CONSTRAINT outbox_status_valid CHECK (status IN ('pending', 'published'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE inbox_events (
    event_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    event_type VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    processed_at DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
