-- Older schema cannot represent non-browser clients. Export them before rollback.
DELETE FROM oauth_clients WHERE kind <> 'authorization_code';

ALTER TABLE oauth_clients
    DROP CHECK oauth_clients_kind_redirect,
    DROP COLUMN kind,
    ADD CONSTRAINT oauth_clients_redirect_uri_nonempty CHECK (redirect_uri <> '');
