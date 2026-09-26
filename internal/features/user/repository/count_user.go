package repository

import "context"

func (ur *UserRepository) CountByAge(
	ctx context.Context,
	minAge *int,
	maxAge *int,
) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM myapp.users
	WHERE ($1::int IS NULL OR users.age >= $1)
	AND ($2::int IS NULL OR users.age <= $2)
`
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
