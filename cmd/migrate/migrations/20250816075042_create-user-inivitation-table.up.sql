CREATE TABLE IF NOT EXISTS "User_Invitations" (
	token bytea PRIMARY KEY,
    user_id bigint NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);