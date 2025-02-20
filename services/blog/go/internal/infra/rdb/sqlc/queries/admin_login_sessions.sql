-- name: ReadAdminLoginSessionByID :one
SELECT id, expired_at FROM admin_login_sessions WHERE id = $1 AND expired_at >= $2;

-- name: CreateAdminLoginSession :exec
INSERT INTO admin_login_sessions(id, created_at, expired_at) VALUES ($1, $2, $3);
