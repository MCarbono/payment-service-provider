package uow

import (
	"database/sql"
	"payment-service-provider/application/repository"
	"payment-service-provider/application/uow"
	infraRepo "payment-service-provider/infra/repository"
	"payment-service-provider/tests/fake"
)

type repositories struct {
	tx *sql.Tx
}

func NewRepositories(tx *sql.Tx) uow.Repositories {
	if f := fake.GetFakeRepositories(tx); f != nil {
		return f
	}
	return &repositories{tx}
}

func (r *repositories) Payable() repository.PayableRepository {
	return infraRepo.NewPayableRepositoryWithTx(r.tx)
}

func (r *repositories) Transaction() repository.TransationRepository {
	return infraRepo.NewTransactionRepositoryWithTX(r.tx)
}
