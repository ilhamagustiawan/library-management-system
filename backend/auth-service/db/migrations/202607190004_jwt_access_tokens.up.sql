ALTER TABLE oauth_tokens
    MODIFY COLUMN access_token VARCHAR(2048) CHARACTER SET ascii COLLATE ascii_bin NULL;
