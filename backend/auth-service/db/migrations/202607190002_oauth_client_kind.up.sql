ALTER TABLE oauth_clients
    DROP CHECK oauth_clients_redirect_uri_nonempty,
    ADD COLUMN kind ENUM('authorization_code', 'client_credentials', 'resource_server')
        NOT NULL DEFAULT 'authorization_code' AFTER name,
    ADD CONSTRAINT oauth_clients_kind_redirect CHECK (
        (kind = 'authorization_code' AND redirect_uri <> '')
        OR (kind IN ('client_credentials', 'resource_server') AND redirect_uri = '')
    );
