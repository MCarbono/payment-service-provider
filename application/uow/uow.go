package uow

import (
	"context"
	"payment-service-provider/application/repository"
)

type UnitOfWorkFn func(ctx context.Context, repositories Repositories) error

type UnitOfWork interface {
	RunTx(ctx context.Context, fn UnitOfWorkFn) error
}

type Repositories interface {
	Payable() repository.PayableRepository
	Transaction() repository.TransationRepository
}
