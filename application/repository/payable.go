package repository

import (
	"context"
	"payment-service-provider/domain/entity"
)

type PayableRepository interface {
	Save(ctx context.Context, payable entity.PayableInterface) error
	GetByID(ctx context.Context, ID string) (entity.PayableInterface, error)
}
