ALTER TABLE oauth_tokens
    MODIFY COLUMN access_token VARCHAR(255) CHARACTER SET ascii COLLATE ascii_bin NULL;
