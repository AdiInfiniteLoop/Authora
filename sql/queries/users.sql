-- name: CreateUser :one
INSERT INTO users (
                   id,
                   name,
                   userName,
                   email,
                   password
)   VALUES (
            $1, $2, $3, $4, $5
           ) RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;