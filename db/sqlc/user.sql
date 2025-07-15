
-- User CRUD

-- name: InsertUser :execresult
INSERT INTO _user (
  email,
  password
) VALUES (?, ?);

-- name: InsertUserIfNotExist :execresult
INSERT INTO _user (email, password)
SELECT ?, ?
WHERE NOT EXISTS (
    SELECT 1 FROM _user WHERE _user.email = ?
);

-- name: UpdateUser :exec
UPDATE _user
SET email = ?, password = ?, is_verified = ?, is_deleted = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM _user
WHERE id = ?;

-- name: GetUserById :one
SELECT * FROM _user
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM _user
WHERE email = ?;

-- name: ListUser :many
SELECT * FROM _user;

