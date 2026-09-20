package repository

import (
	"context"
	"errors"
	"log/slog"
	"myapp/internal/core/domains/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(db *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{db, logger}
}

func (ur *UserRepository) FindAll(ctx context.Context, minAge int, maxAge int, limit int, offset int) ([]user.User, error) {
	query := "SELECT users.id, users.name, users.age from myapp.users where users.age BETWEEN $1 AND $2 LIMIT $3 OFFSET $4"
	rows, err := ur.db.Query(ctx, query, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, DatabaseError
	}
	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User

		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.Age,
		); err != nil {
			return nil, DatabaseError
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, DatabaseError
	}

	return users, nil
}

func (ur *UserRepository) FindByID(ctx context.Context, id int) (*user.User, error) {
	query := "SELECT users.id, users.name, users.age from myapp.users where users.id = $1"

	var u user.User

	err := ur.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Age,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return &u, UserNotFound
	}
	if err != nil {
		return &u, DatabaseError
	}
	return &u, nil
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id int) error {
	query := "DELETE FROM myapp.users WHERE id = $1"
	commandTag, err := ur.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return DatabaseError
	}
	if commandTag.RowsAffected() == 0 {
		return UserNotFound
	}
	return nil
}

func (ur *UserRepository) Update(ctx context.Context, id int, u *user.User) (*user.User, error) {
	ur.logger.Info("updating user")
	query := "UPDATE myapp.users SET name = $1, age = $2 WHERE id = $3"
	queryResult, err := ur.db.Exec(
		ctx,
		query,
		u.Name,
		u.Age,
		id,
	)
	if err != nil {
		ur.logger.Error(err.Error())
		return &user.User{}, DatabaseError
	}
	if queryResult.RowsAffected() == 0 {
		ur.logger.Warn("user not found")
		return &user.User{}, UserNotFound
	}
	ur.logger.Info("user updated")
	return ur.FindByID(ctx, id)
}

func (ur *UserRepository) Save(ctx context.Context, u *user.User) (*user.User, error) {
	query := "INSERT INTO myapp.users (name, age) VALUES ($1, $2) RETURNING id, name, age"

	err := ur.db.QueryRow(
		ctx,
		query,
		u.Name,
		u.Age,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Age,
	)
	if err != nil {
		return u, DatabaseError
	}
	return u, nil
}

func (ur *UserRepository) CountByAge(ctx context.Context, minAge int, maxAge int) (int, error) {
	query := "SELECT COUNT(*) FROM myapp.users WHERE age BETWEEN $1 AND $2"
	rows, err := ur.db.Query(ctx, query, minAge, maxAge)
	if err != nil {
		return 0, DatabaseError
	}
	defer rows.Close()
	rows.Next()
	var count int
	err = rows.Scan(&count)
	if err != nil {
		return 0, DatabaseError
	}
	return count, nil
}
