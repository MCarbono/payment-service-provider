package uow

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service-provider/application/uow"
	"payment-service-provider/infra/repository"
)

type UowImpl struct {
	db *sql.DB
}

func NewUow(db *sql.DB) *UowImpl {
	return &UowImpl{db: db}
}

func (u *UowImpl) RunTx(ctx context.Context, fn uow.UnitOfWorkFn) error {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error trying to open tx: %w", err)
	}
	err = fn(ctx, repository.NewRepositories(tx))
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("commit error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
	}
	return nil
}
