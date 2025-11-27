-- name: InsertRefreshToken :exec
INSERT INTO refresh_tokens (
  hash_token,
  user_id,
  expires_at
) VALUES ( $1, $2, $3 );

-- name: GetRefreshTokensByUserID :many
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE hash_token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE hash_token = $1;
