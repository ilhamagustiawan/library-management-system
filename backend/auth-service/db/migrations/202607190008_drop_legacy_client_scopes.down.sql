ALTER TABLE oauth_clients
    ADD COLUMN allowed_scopes TEXT NOT NULL AFTER redirect_uri;

UPDATE oauth_clients
SET allowed_scopes = COALESCE((
    SELECT GROUP_CONCAT(scope_code ORDER BY scope_code SEPARATOR ' ')
    FROM oauth_client_scopes
    WHERE oauth_client_scopes.client_id = oauth_clients.id
), '');
