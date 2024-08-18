package uow

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service-provider/application/uow"
	"payment-service-provider/infra/tracing"
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
	_, span := tracing.Tracer.Start(ctx, "UOW - execution")
	defer span.End()
	err = fn(ctx, NewRepositories(tx))
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	_, span2 := tracing.Tracer.Start(ctx, "UOW - commiting")
	defer span2.End()
	err = tx.Commit()
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return fmt.Errorf("commit error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
	}
	return nil
}
