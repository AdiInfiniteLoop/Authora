// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
                   id,
                   name,
                   userName,
                   email,
                   password
)   VALUES (
            $1, $2, $3, $4, $5
           ) RETURNING id, name, username, email, password, created_at, updated_at
`

type CreateUserParams struct {
	ID       uuid.UUID
	Name     string
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Username,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, name, username, email, password, created_at, updated_at FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
