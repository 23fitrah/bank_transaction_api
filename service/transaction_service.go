package service

import (
	"context"
	"fmt"
	"log"
	"transaction_api/dto/transaction"
	model "transaction_api/model/transaction"
	"transaction_api/repository"
)

type TransactionService struct {
	Repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		Repo: repo,
	}
}

func (s *TransactionService) SaveTransactionService(ctx context.Context, data transaction.RequestField) (int64, error) {

	dataTransaction := model.Transaction{
		ROWID_SENDER:      data.ROWID_SENDER,
		ROWID_BENEFICIARY: data.ROWID_BENEFICIARY,
		TRANSACTION_DATE:  data.TRANSACTION_DATE,
		DEBET_ACCOUNT:     data.DEBET_ACCOUNT,
		DEBET_NAME:        data.DEBET_NAME,
		DEBET_CURR:        data.DEBET_CURR,
		CREDIT_ACCOUNT:    data.CREDIT_ACCOUNT,
		CREDIT_NAME:       data.CREDIT_NAME,
		CREDIT_CURR:       data.CREDIT_CURR,
		AMOUNT:            data.AMOUNT,
		REMARK:            data.REMARK,
		REFERENCE_NUMBER:  data.REFERENCE_NUMBER,
	}
	insertedID, err := s.Repo.SaveTransactionRepository(
		ctx,
		dataTransaction,
	)
	if err != nil {
		log.Println(err)

		return 0, fmt.Errorf("failed to insert user fraud: %w", err)
	}

	return insertedID, nil
}

func (s *TransactionService) CheckIdOriginalExistsService(ctx context.Context, idOriginal string) (int64, error) {

	nextID, err := s.Repo.CheckIdOriginalExistsRepository(ctx, idOriginal)
	if err != nil {

		return 0, fmt.Errorf("failed to get next ID: %w", err)
	}

	return int64(nextID), nil
}

func (s *TransactionService) GetAllTransactionsService(ctx context.Context, limit, offset int, refNo, dateFrom, dateTo string) ([]*model.Transaction, int64, error) {
	list, total, err := s.Repo.GetAllTransactionRepository(ctx, limit, offset, refNo, dateFrom, dateTo)
	if err != nil {

		return nil, 0, fmt.Errorf("service get all transactions: %w", err)
	}
	return list, total, nil
}

func (s *TransactionService) GetDetailTransactionService(ctx context.Context, rowId string) (*model.Transaction, error) {
	uf, err := s.Repo.GetDetailTransactionRepository(ctx, rowId)
	if err != nil {
		return nil, fmt.Errorf("get detail transaction: %w", err)
	}
	return uf, nil
}
