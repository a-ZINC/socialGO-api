CREATE TABLE IF NOT EXISTS user_invitations (
	token bytea PRIMARY KEY,
    user_id bigint NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);