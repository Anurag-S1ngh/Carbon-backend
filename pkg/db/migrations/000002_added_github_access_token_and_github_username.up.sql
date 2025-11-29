ALTER TABLE users
  ADD COLUMN github_username VARCHAR(255) UNIQUE,
  ADD COLUMN github_access_token TEXT;

