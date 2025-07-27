CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_comment_contents ON "Comment" USING GIN (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_post_title ON "Post" USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_post_tags ON "Post" USING GIN (tags gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_comment_post_id ON "Comment" (post_id);
CREATE INDEX IF NOT EXISTS idx_post_user_id ON "Post" (user_id);
CREATE INDEX IF NOT EXISTS idx_user_name ON "User" (name)