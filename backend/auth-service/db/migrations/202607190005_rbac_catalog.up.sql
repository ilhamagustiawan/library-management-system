CREATE TABLE roles (
    code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    description VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE scopes (
    code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
    audience VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    description VARCHAR(255) NOT NULL,
    KEY scopes_audience_idx (audience)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE role_scopes (
    role_code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    scope_code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    PRIMARY KEY (role_code, scope_code),
    CONSTRAINT role_scopes_role_fk FOREIGN KEY (role_code) REFERENCES roles(code),
    CONSTRAINT role_scopes_scope_fk FOREIGN KEY (scope_code) REFERENCES scopes(code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO roles (code, description) VALUES
    ('member', 'Library member'),
    ('admin', 'Library administrator');

INSERT INTO scopes (code, audience, description) VALUES
    ('loans:borrow:self', 'library-api', 'Borrow for token subject'),
    ('loans:return:self', 'library-api', 'Return token subject loan'),
    ('transactions:read:self', 'library-api', 'Read token subject transaction history'),
    ('books:read', 'library-api', 'Browse catalog and availability'),
    ('transactions:read:any', 'library-api', 'Read all member transactions'),
    ('loans:return:any', 'library-api', 'Return any member loan'),
    ('fines:manage', 'library-api', 'Manage fine records'),
    ('books:manage', 'library-api', 'Manage catalog and inventory'),
    ('identities:create', 'auth-service', 'Create member identities'),
    ('book-stock:read', 'book-service', 'Read book stock'),
    ('book-stock:reserve', 'book-service', 'Reserve book stock'),
    ('book-stock:release', 'book-service', 'Release book stock');

INSERT INTO role_scopes (role_code, scope_code) VALUES
    ('member', 'loans:borrow:self'),
    ('member', 'loans:return:self'),
    ('member', 'transactions:read:self'),
    ('member', 'books:read'),
    ('admin', 'transactions:read:any'),
    ('admin', 'loans:return:any'),
    ('admin', 'fines:manage'),
    ('admin', 'books:manage');

ALTER TABLE users
    ADD COLUMN role_code VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin
        NOT NULL DEFAULT 'member' AFTER password_hash,
    ADD CONSTRAINT users_role_fk FOREIGN KEY (role_code) REFERENCES roles(code);
