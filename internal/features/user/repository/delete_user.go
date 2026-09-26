package repository

import "context"

func (ur *UserRepository) DeleteByID(
	ctx context.Context,
	id int) error {
	query := `
	DELETE FROM myapp.users 
    WHERE id = $1`

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
