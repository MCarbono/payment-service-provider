package fake

import (
	"context"
	"database/sql"
	"errors"
	"payment-service-provider/application/repository"
	"payment-service-provider/domain/entity"
	infraRepo "payment-service-provider/infra/repository"
)

type fakeRepositories interface {
	setTx(tx *sql.Tx)
}

var fakeRepositoriesImpl fakeRepositories

func GetFakeRepositories(tx *sql.Tx) fakeRepositories {
	if fakeRepositoriesImpl != nil {
		fakeRepositoriesImpl.setTx(tx)
		return fakeRepositoriesImpl
	}
	return nil
}

func NewFakeRepositoriesSavePayableError() {
	fakeRepositoriesImpl = &fakeRepositoriesSavePayableError{}
}

type fakeRepositoriesSavePayableError struct {
	tx *sql.Tx
}

func (f *fakeRepositoriesSavePayableError) setTx(tx *sql.Tx) {
	f.tx = tx
}

func (f *fakeRepositoriesSavePayableError) Payable() repository.PayableRepository {
	return &fakePayableRepositorySaveError{}
}

func (f *fakeRepositoriesSavePayableError) Transaction() repository.TransationRepository {
	return infraRepo.NewTransactionRepositoryWithTX(f.tx)
}

type fakePayableRepositorySaveError struct {
}

func (f *fakePayableRepositorySaveError) Save(ctx context.Context, payable *entity.Payable) error {
	return errors.New("internal server error")
}

func (f *fakePayableRepositorySaveError) GetByID(ctx context.Context, ID string) (*entity.Payable, error) {
	return nil, nil
}

func NewFakeRepositoriesSaveTransactionError() {
	fakeRepositoriesImpl = &fakeRepositoriesSaveTransactionError{}
}

func DestroyFakeRepositoriesImpl() {
	fakeRepositoriesImpl = nil
}

type fakeRepositoriesSaveTransactionError struct {
	tx *sql.Tx
}

func (f *fakeRepositoriesSaveTransactionError) setTx(tx *sql.Tx) {
	f.tx = tx
}

func (f *fakeRepositoriesSaveTransactionError) Payable() repository.PayableRepository {
	return infraRepo.NewPayableRepositoryWithTx(f.tx)
}

func (f *fakeRepositoriesSaveTransactionError) Transaction() repository.TransationRepository {
	return &fakeTransactionRepositoryError{}
}

type fakeTransactionRepositoryError struct {
}

func (f *fakeTransactionRepositoryError) Save(ctx context.Context, transaction *entity.Transaction) error {
	return errors.New("internal server error")
}

func (f *fakeTransactionRepositoryError) GetByID(ctx context.Context, ID string) (*entity.Transaction, error) {
	return nil, nil
}
