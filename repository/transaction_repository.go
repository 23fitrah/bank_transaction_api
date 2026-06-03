package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
	model "transaction_api/model/transaction"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) SaveTransactionRepository(ctx context.Context, data model.Transaction) (int64, error) {

	const query = `
		INSERT INTO 
			[t_transactions] (
			ROWID_SENDER,
			ROWID_BENEFICIARY,
			TRANSACTION_DATE,
			DEBET_ACCOUNT,
			DEBET_NAME,
			DEBET_CURR,
			CREDIT_ACCOUNT,
			CREDIT_NAME,
			CREDIT_CURR,
			AMOUNT,
			REMARK,
			REFERENCE_NUMBER,
			STATUS_CODE,
			STATUS_DESC
		)
		OUTPUT 
			INSERTED.ROWID_TRX
		VALUES (
			@p1,            
			@p2,            
			GETDATE(),            
			@p3,            
			@p4,            
			@p5,            
			@p6,            
			@p7,            
			@p8,            
			@p9,            
			@p10,           
			@p11,           
			@p12,
			@p13            
        );
    `

	t := time.Now()
	formatted := t.Format("20060102150405")
	refno := "0" + formatted

	res, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", data.ROWID_SENDER),
		sql.Named("p2", data.ROWID_BENEFICIARY),
		sql.Named("p3", data.DEBET_ACCOUNT),
		sql.Named("p4", data.DEBET_NAME),
		sql.Named("p5", data.DEBET_CURR),
		sql.Named("p6", data.CREDIT_ACCOUNT),
		sql.Named("p7", data.CREDIT_NAME),
		sql.Named("p8", data.CREDIT_CURR),
		sql.Named("p9", data.AMOUNT),
		sql.Named("p10", data.REMARK),
		sql.Named("p11", refno),
		sql.Named("p12", "0000"),
		sql.Named("p13", "Ready to process"))

	if err != nil {
		return 0, fmt.Errorf("error inserting into t_transaction: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error inserting rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (r *TransactionRepository) CheckIdOriginalExistsRepository(ctx context.Context, idOriginal string) (int, error) {

	var count int

	idOriginalInt, err := strconv.ParseInt(idOriginal, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ROWID_TRX format: %w", err)
	}
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM t_transactions WHERE ROWID_TRX = ?", idOriginalInt).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check ROWID_TRX: %w", err)
	}

	if count > 0 {
		return 1, nil
	}
	return 0, nil
}

func (r *TransactionRepository) GetAllTransactionRepository(ctx context.Context, limit, offset int, refNo, dateFrom, dateTo string) ([]*model.Transaction, int64, error) {

	query := `
		SELECT 
			ROWID_TRX,
			ROWID_SENDER,
			ROWID_BENEFICIARY,
			REFERENCE_NUMBER,
			STATUS_CODE,
			STATUS_DESC,
			TRANSACTION_DATE,
			DEBET_ACCOUNT,
			DEBET_NAME,
			DEBET_CURR,
			CREDIT_ACCOUNT,
			CREDIT_NAME,
			CREDIT_CURR,
			AMOUNT,
			REMARK
		FROM 
			t_transactions WITH (NOLOCK)
		WHERE 
			1=1 `

	args := []interface{}{}
	paramIdx := 1

	if refNo != "" {
		query = query + " AND REFERENCE_NUMBER = @p" + strconv.Itoa(paramIdx)
		args = append(args, refNo)
		paramIdx++
	}

	if dateFrom != "" && dateTo != "" {
		query = query + " AND TRANSACTION_DATE BETWEEN @p" + strconv.Itoa(paramIdx) + " AND @p" + strconv.Itoa(paramIdx+1)
		args = append(args, dateFrom, dateTo)
		paramIdx += 2
	}

	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM (`+query+`) AS tb`, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transaction: %w", err)
	}

	query += " ORDER BY ROWID_TRX DESC"

	if limit != -1 {
		query += fmt.Sprintf("  OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", paramIdx, paramIdx+1)
	}
	args = append(args, offset)
	args = append(args, limit)

	sqlRows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, 0, fmt.Errorf("query get all transaction failed: %w", err)
	}

	defer sqlRows.Close()

	results := []*model.Transaction{}
	for sqlRows.Next() {
		var trx model.Transaction
		err := sqlRows.Scan(
			&trx.ROWID_TRX,
			&trx.ROWID_SENDER,
			&trx.ROWID_BENEFICIARY,
			&trx.REFERENCE_NUMBER,
			&trx.STATUS_CODE,
			&trx.STATUS_DESC,
			&trx.TRANSACTION_DATE,
			&trx.DEBET_ACCOUNT,
			&trx.DEBET_NAME,
			&trx.DEBET_CURR,
			&trx.CREDIT_ACCOUNT,
			&trx.CREDIT_NAME,
			&trx.CREDIT_CURR,
			&trx.AMOUNT,
			&trx.REMARK)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction row failed: %w", err)
		}
		results = append(results, &trx)
	}

	return results, total, nil
}

func (r *TransactionRepository) GetDetailTransactionRepository(ctx context.Context, rowId string) (*model.Transaction, error) {
	const query = `
		SELECT TOP 1 
			ROWID_TRX, 
			ROWID_SENDER, 
			ROWID_BENEFICIARY, 
			REFERENCE_NUMBER,
			STATUS_CODE,
			STATUS_DESC,
			TRANSACTION_DATE, 
			DEBET_ACCOUNT, 
			DEBET_NAME, 
			DEBET_CURR, 
			CREDIT_ACCOUNT, 
			CREDIT_NAME, 
			CREDIT_CURR, 
			AMOUNT, 
			REMARK
		FROM 
			t_transactions 
		WHERE 
			ROWID_TRX = @p1 `
	var uf model.Transaction

	err := r.db.QueryRowContext(ctx, query, rowId).Scan(
		&uf.ROWID_TRX,
		&uf.ROWID_SENDER,
		&uf.ROWID_BENEFICIARY,
		&uf.REFERENCE_NUMBER,
		&uf.STATUS_CODE,
		&uf.STATUS_DESC,
		&uf.TRANSACTION_DATE,
		&uf.DEBET_ACCOUNT,
		&uf.DEBET_NAME,
		&uf.DEBET_CURR,
		&uf.CREDIT_ACCOUNT,
		&uf.CREDIT_NAME,
		&uf.CREDIT_CURR,
		&uf.AMOUNT,
		&uf.REMARK)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("query get detail transaction: %w", err)
	}

	return &uf, nil
}
