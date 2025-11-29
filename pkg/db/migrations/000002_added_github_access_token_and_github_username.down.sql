ALTER TABLE users
  DROP COLUMN IF EXISTS github_access_token,
  DROP COLUMN IF EXISTS github_username;

