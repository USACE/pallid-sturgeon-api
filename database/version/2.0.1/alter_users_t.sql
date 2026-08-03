ALTER TABLE users_t ADD token_access_id VARCHAR2(40);
ALTER TABLE users_t ADD token_secret_id VARCHAR2(40);
ALTER TABLE users_t ADD token_expiration DATE;