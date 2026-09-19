package repository

import (
	"context"
	"log/slog"
	"myapp/internal/core/domains/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	ctx    context.Context
	logger *slog.Logger
}

func New(db *pgxpool.Pool, ctx context.Context, logger *slog.Logger) *UserRepository {
	return &UserRepository{db, ctx, logger}
}

func (ur *UserRepository) FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, error) {
	query := "SELECT users.id, users.name, users.age from myapp.users"
	rows, err := ur.db.Query(ur.ctx, query, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, nil
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

func (ur *UserRepository) FindByID(id int) (user.User, error) {
	query := "SELECT users.id, users.name, users.age from myapp.users where users.id = $1"

	var u user.User

	err := ur.db.QueryRow(
		ur.ctx,
		query,
		id,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Age,
	)
	if err != nil {
		return u, UserNotFound
	}
	return u, nil
}

func (ur *UserRepository) DeleteByID(id int) error {
	query := "DELETE FROM myapp.users WHERE id = $1"
	commandTag, err := ur.db.Exec(
		ur.ctx,
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

func (ur *UserRepository) Update(id int, name string, age int) (user.User, error) {
	ur.logger.Info("updating user")
	query := "UPDATE myapp.users SET name = $1, age = $2 WHERE id = $3"
	queryResult, err := ur.db.Exec(
		ur.ctx,
		query,
		name,
		age,
		id,
	)
	ur.logger.Info("user updated")
	if err != nil {
		ur.logger.Error(err.Error())
		return user.User{}, DatabaseError
	}
	if queryResult.RowsAffected() == 0 {
		ur.logger.Warn("user not found")
		return user.User{}, UserNotFound
	}
	ur.logger.Info("user updated")
	return ur.FindByID(id)
}

func (ur *UserRepository) Create(name string, age int) (user.User, error) {
	query := "INSERT INTO myapp.users (name, age) VALUES ($1, $2) RETURNING id, name, age"

	var u user.User
	err := ur.db.QueryRow(
		ur.ctx,
		query,
		name,
		age,
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

func (ur *UserRepository) CountByAge(minAge int, maxAge int) (int, error) {
	query := "SELECT COUNT(*) FROM myapp.users WHERE age BETWEEN $1 AND $2"
	rows, err := ur.db.Query(ur.ctx, query, minAge, maxAge)
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
