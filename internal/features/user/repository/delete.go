package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
		DELETE FROM myapp.users
		WHERE id = $1
	`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return postgres.MapError(err)
	}

	if tag.RowsAffected() == 0 {
		return errs.NotFound("user not found", nil)
	}

	return nil
}
