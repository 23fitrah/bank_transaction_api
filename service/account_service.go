package service

import (
	"context"
	"fmt"
	dto "transaction_api/dto/account"
	"transaction_api/model/account"
	model "transaction_api/model/account"
	"transaction_api/repository"
)

type AccountService struct {
	Repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{
		Repo: repo,
	}
}

func (s *AccountService) GetInquiryAccountService(ctx context.Context, accountNo string) (*account.Account, error) {
	uf, err := s.Repo.GetInquiryAccountRepository(ctx, accountNo)
	if err != nil {
		return nil, fmt.Errorf("get inquiry account: %w", err)
	}
	return uf, nil
}

func (s *AccountService) InsertAccountService(ctx context.Context, data dto.AccountRequestField) (int64, error) {
	dataMaster := model.Account{
		ACCOUNT_NO: data.ACCOUNT_NO,
		NAME:       data.NAME,
		CURRENCY:   data.CURRENCY,
		BALANCE:    data.BALANCE,
		STATUS:     data.STATUS,
		EMAIL:      data.EMAIL,
		ADDRESS:    data.ADDRESS,
		COUNTRY:    data.COUNTRY,
	}

	insertID, err := s.Repo.InsertAccountRepository(ctx, dataMaster)
	if err != nil {
		return 0, fmt.Errorf("insert account: %w", err)
	}
	return insertID, nil
}

func (s *AccountService) UpdateAccountService(ctx context.Context, data dto.AccountRequestField) (int64, error) {
	dataMaster := model.Account{
		ID_ACCOUNT: data.ID_ACCOUNT,
		ACCOUNT_NO: data.ACCOUNT_NO,
		NAME:       data.NAME,
		CURRENCY:   data.CURRENCY,
		BALANCE:    data.BALANCE,
		STATUS:     data.STATUS,
		EMAIL:      data.EMAIL,
		ADDRESS:    data.ADDRESS,
		COUNTRY:    data.COUNTRY,
	}
	updateID, err := s.Repo.UpdateAccountRepository(ctx, dataMaster)
	if err != nil {
		return 0, fmt.Errorf("update account: %w", err)
	}
	return updateID, nil
}

func (s *AccountService) DeleteAccountService(ctx context.Context, id int64) (int64, error) {
	deleteID, err := s.Repo.DeleteAccountRepository(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("delete account: %w", err)
	}
	return deleteID, nil
}

func (s *AccountService) CheckAccountIDExistsService(ctx context.Context, id int64) (int, error) {
	exist, err := s.Repo.CheckAccountIDExistsRepository(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("get detail account: %w", err)
	}
	return exist, nil
}

func (s *AccountService) CheckAccountNoExistsService(ctx context.Context, accountNo string) (int, error) {
	exists, err := s.Repo.CheckAccountNoExistsRepository(ctx, accountNo)
	if err != nil {
		return 0, fmt.Errorf("check account no exists: %w", err)
	}
	return exists, nil
}

func (s *AccountService) GetDetailAccountService(ctx context.Context, id int64) (*account.Account, error) {
	uf, err := s.Repo.GetDetailAccountRepository(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get detail account: %w", err)
	}
	return uf, nil
}

func (s *AccountService) GetAllAccountsService(ctx context.Context, limit, offset int, accountNo, name string) ([]*account.Account, int64, error) {
	list, total, err := s.Repo.GetAllAccountsRepository(ctx, limit, offset, accountNo, name)
	if err != nil {
		return nil, 0, fmt.Errorf("get all accounts: %w", err)
	}

	return list, total, nil
}
