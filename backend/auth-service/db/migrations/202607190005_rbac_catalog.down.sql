ALTER TABLE users
    DROP FOREIGN KEY users_role_fk,
    DROP COLUMN role_code;

DROP TABLE role_scopes;
DROP TABLE scopes;
DROP TABLE roles;
