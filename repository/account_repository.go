package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"transaction_api/model/account"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetInquiryAccountRepository(ctx context.Context, accountNo string) (*account.Account, error) {
	const query = `
		SELECT TOP 1 
			ACCOUNT_NO, 
			NAME, 
			BALANCE, 
			CURRENCY,
			STATUS 
		FROM 
			m_account 
		WHERE 
			ACCOUNT_NO = @p1 `
	var uf account.Account

	err := r.db.QueryRowContext(ctx, query, accountNo).Scan(
		&uf.ACCOUNT_NO,
		&uf.NAME,
		&uf.BALANCE,
		&uf.CURRENCY,
		&uf.STATUS,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("query get detail account our: %w", err)
	}

	return &uf, nil
}

func (r *AccountRepository) UpdateAccountRepository(ctx context.Context, data account.Account) (int64, error) {
	const query = `
		UPDATE [m_account]
		SET
			NAME 		= @p2,
			ACCOUNT_NO	= @p3,
			ADDRESS 	= @p4,
			CURRENCY 	= @p5,
			STATUS 		= @p6,
			COUNTRY	 	= @p7,
			EMAIL 		= @p8,
			BALANCE 	= @p9
		WHERE ID_ACCOUNT= @p1;
	`

	res, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", data.ID_ACCOUNT),
		sql.Named("p2", data.NAME),
		sql.Named("p3", data.ACCOUNT_NO),
		sql.Named("p4", data.ADDRESS),
		sql.Named("p5", data.CURRENCY),
		sql.Named("p6", strings.ToUpper(data.STATUS)),
		sql.Named("p7", data.COUNTRY),
		sql.Named("p8", data.EMAIL),
		sql.Named("p9", data.BALANCE))
	if err != nil {
		return 0, fmt.Errorf("error updating account: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error updateing rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (r *AccountRepository) InsertAccountRepository(ctx context.Context, data account.Account) (int64, error) {
	const query = `
		INSERT INTO [m_account] (
			NAME,
			ACCOUNT_NO,
			ADDRESS,
			COUNTRY,
			EMAIL,
			BALANCE,
			CURRENCY,
			STATUS
		)
		OUTPUT INSERTED.ID_ACCOUNT
		VALUES (
			@p1,            
			@p2,	
			@p3,
			@p4,
			@p5,
			@p6,
			@p7,
			@p8
		);
	`

	res, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", data.NAME),
		sql.Named("p2", data.ACCOUNT_NO),
		sql.Named("p3", data.ADDRESS),
		sql.Named("p4", data.COUNTRY),
		sql.Named("p5", data.EMAIL),
		sql.Named("p6", data.BALANCE),
		sql.Named("p7", data.CURRENCY),
		sql.Named("p8", strings.ToUpper(data.STATUS)))

	if err != nil {
		return 0, fmt.Errorf("error inserting into account: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error inserting rows affected: %w", err)
	}

	return rowsAffected, nil

}

func (r *AccountRepository) DeleteAccountRepository(ctx context.Context, id int64) (int64, error) {
	const query = `
		DELETE FROM [m_account]
		OUTPUT DELETED.ID_ACCOUNT
		WHERE ID_ACCOUNT = @p1;
	`

	res, err := r.db.ExecContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return 0, fmt.Errorf("error deleting account: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error deleting rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (r *AccountRepository) CheckAccountIDExistsRepository(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT COUNT(1)
		FROM [m_account]
		WHERE ID_ACCOUNT = @p1
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to check ID account: %w", err)
	}
	if count > 0 {
		return 1, nil
	}
	return 0, nil
}

func (r *AccountRepository) CheckAccountNoExistsRepository(ctx context.Context, accountNo string) (int, error) {
	const query = `
		SELECT COUNT(1)	
		FROM [m_account]
		WHERE ACCOUNT_NO = @p1
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query, sql.Named("p1", accountNo)).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to check account number: %w", err)
	}
	if count > 0 {
		return 1, nil
	}
	return 0, nil
}

func (r *AccountRepository) GetDetailAccountRepository(ctx context.Context, id int64) (*account.Account, error) {
	const query = `
		SELECT TOP 1
			ID_ACCOUNT,
			NAME,
			ACCOUNT_NO,
			ADDRESS,
			CURRENCY,
			STATUS,
			COUNTRY,
			EMAIL,
			BALANCE
		FROM [m_account]
		WHERE ID_ACCOUNT = @p1
	`
	var uf account.Account
	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).Scan(
		&uf.ID_ACCOUNT,
		&uf.NAME,
		&uf.ACCOUNT_NO,
		&uf.ADDRESS,
		&uf.CURRENCY,
		&uf.STATUS,
		&uf.COUNTRY,
		&uf.EMAIL,
		&uf.BALANCE,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("query get detail account: %w", err)
	}
	return &uf, nil
}

func (r *AccountRepository) GetAllAccountsRepository(ctx context.Context, limit, offset int, acct_no, name string) ([]*account.Account, int64, error) {
	query := `
		SELECT
			ID_ACCOUNT,
			NAME,
			ACCOUNT_NO,
			ADDRESS,
			CURRENCY,
			STATUS,
			COUNTRY,
			EMAIL,
			BALANCE
		FROM [m_account]  WHERE 1=1
	`
	var conditions []string
	var args []interface{}
	paramIdx := 1

	if acct_no != "" {
		conditions = append(conditions, fmt.Sprintf("ACCOUNT_NO LIKE @p%d", paramIdx))
		args = append(args, "%"+acct_no+"%")
		paramIdx++
	}
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("NAME LIKE @p%d", paramIdx))
		args = append(args, "%"+name+"%")
		paramIdx++
	}
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM (`+query+`) AS tb`, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count account: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY ID_ACCOUNT OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", paramIdx, paramIdx+1)
	args = append(args, offset)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query get all accounts: %w", err)
	}
	defer rows.Close()

	var results []*account.Account
	for rows.Next() {
		var acc account.Account
		err := rows.Scan(
			&acc.ID_ACCOUNT,
			&acc.NAME,
			&acc.ACCOUNT_NO,
			&acc.ADDRESS,
			&acc.CURRENCY,
			&acc.STATUS,
			&acc.COUNTRY,
			&acc.EMAIL,
			&acc.BALANCE,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan account row failed: %w", err)
		}
		results = append(results, &acc)
	}

	return results, total, nil
}
