CREATE TABLE books (
    id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    isbn VARCHAR(13) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT NULL,
    publication_year SMALLINT UNSIGNED NULL,
    total_copies INT UNSIGNED NOT NULL,
    available_copies INT UNSIGNED NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    archived_at DATETIME(6) NULL,
    UNIQUE KEY books_isbn_uidx (isbn),
    KEY books_title_idx (title),
    KEY books_author_idx (author),
    KEY books_archived_at_idx (archived_at),
    CONSTRAINT books_copy_counts_valid CHECK (
        available_copies <= total_copies
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE book_reservations (
    transaction_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    book_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    created_at DATETIME(6) NOT NULL,
    released_at DATETIME(6) NULL,
    KEY book_reservations_book_status_idx (book_id, status),
    CONSTRAINT book_reservations_book_fk FOREIGN KEY (book_id)
        REFERENCES books(id) ON DELETE RESTRICT,
    CONSTRAINT book_reservations_status_valid CHECK (status IN ('active', 'released')),
    CONSTRAINT book_reservations_release_valid CHECK (
        (status = 'active' AND released_at IS NULL)
        OR (status = 'released' AND released_at IS NOT NULL)
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
