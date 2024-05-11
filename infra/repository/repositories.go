package repository

import (
	"database/sql"
	"payment-service-provider/application/repository"
	"payment-service-provider/application/uow"
)

type repositories struct {
	tx *sql.Tx
}

func NewRepositories(tx *sql.Tx) uow.Repositories {
	return &repositories{tx}
}

func (r *repositories) Payable() repository.PayableRepository {
	return NewPayableRepositoryWithTx(r.tx)
}

func (r *repositories) Transaction() repository.TransationRepository {
	return NewTransactionRepositoryWithTX(r.tx)
}
