CREATE TABLE oauth_client_scopes (
    client_id VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    scope_code VARCHAR(100) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    PRIMARY KEY (client_id, scope_code),
    CONSTRAINT oauth_client_scopes_client_fk FOREIGN KEY (client_id)
        REFERENCES oauth_clients(id) ON DELETE CASCADE,
    CONSTRAINT oauth_client_scopes_scope_fk FOREIGN KEY (scope_code)
        REFERENCES scopes(code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Translate the only legacy grants. Unrecognized values receive no mapping and fail closed.
INSERT IGNORE INTO oauth_client_scopes (client_id, scope_code)
SELECT id, 'books:read' FROM oauth_clients
WHERE kind = 'authorization_code' AND FIND_IN_SET('library:read', REPLACE(allowed_scopes, ' ', ',')) > 0;

INSERT IGNORE INTO oauth_client_scopes (client_id, scope_code)
SELECT id, 'transactions:read:self' FROM oauth_clients
WHERE kind = 'authorization_code' AND FIND_IN_SET('library:read', REPLACE(allowed_scopes, ' ', ',')) > 0;

INSERT IGNORE INTO oauth_client_scopes (client_id, scope_code)
SELECT id, 'loans:borrow:self' FROM oauth_clients
WHERE kind = 'authorization_code' AND FIND_IN_SET('library:write', REPLACE(allowed_scopes, ' ', ',')) > 0;

INSERT IGNORE INTO oauth_client_scopes (client_id, scope_code)
SELECT id, 'loans:return:self' FROM oauth_clients
WHERE kind = 'authorization_code' AND FIND_IN_SET('library:write', REPLACE(allowed_scopes, ' ', ',')) > 0;
