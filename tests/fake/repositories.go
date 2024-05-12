package fake

import (
	"context"
	"database/sql"
	"errors"
	"payment-service-provider/application/repository"
	"payment-service-provider/application/uow"
	"payment-service-provider/domain/entity"
	infraRepo "payment-service-provider/infra/repository"
)

func GetFakeRepositories(tx *sql.Tx) uow.Repositories {
	if fakeRepositoriesSavePayableErrorSingleton != nil {
		fakeRepositoriesSavePayableErrorSingleton.setTx(tx)
		return fakeRepositoriesSavePayableErrorSingleton
	}
	if fakeRepositoriesTransactionSaveErrorSingleton != nil {
		fakeRepositoriesTransactionSaveErrorSingleton.setTx(tx)
		return fakeRepositoriesTransactionSaveErrorSingleton
	}
	return nil
}

var fakeRepositoriesSavePayableErrorSingleton *fakeRepositoriesSavePayableError

func NewFakeRepositoriesSavePayableError() {
	if fakeRepositoriesSavePayableErrorSingleton == nil {
		fakeRepositoriesSavePayableErrorSingleton = &fakeRepositoriesSavePayableError{}
	}
}

func DestroyFakeRepositoriesSavePayableError() {
	fakeRepositoriesSavePayableErrorSingleton = nil
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

var fakeRepositoriesTransactionSaveErrorSingleton *fakeRepositoriesSaveTransactionError

func NewFakeRepositoriesSaveTransactionError() {
	if fakeRepositoriesTransactionSaveErrorSingleton == nil {
		fakeRepositoriesTransactionSaveErrorSingleton = &fakeRepositoriesSaveTransactionError{}
	}
}

func DestroyFakeRepositoriesTransactionError() {
	fakeRepositoriesTransactionSaveErrorSingleton = nil
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
