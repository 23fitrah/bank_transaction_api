package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"transaction_api/model/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetInquiryUserRepository(ctx context.Context, username string) (*user.User, error) {
	const query = `
		SELECT TOP 1 
			id, 
			name, 
			email, 
			password,
			role,
			token,
			updatedat
		FROM 
			users
		WHERE 
			name = @p1 `
	var uf user.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&uf.Id,
		&uf.Name,
		&uf.Email,
		&uf.Password,
		&uf.Role,
		&uf.Token,
		&uf.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("query get detail user: %w", err)
	}

	return &uf, nil
}

func (r *UserRepository) UpdateUserRepository(ctx context.Context, id int, username, token string) (int64, error) {
	const query = `
		UPDATE users
		SET
			token = @p2,
			updatedat=GETDATE()
		WHERE id = @p1 AND name=@p3;
	`

	res, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", id),
		sql.Named("p2", token),
		sql.Named("p3", username))
	if err != nil {
		return 0, fmt.Errorf("error updating user token: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error updateing user token rows affected: %w", err)
	}

	return rowsAffected, nil
}
