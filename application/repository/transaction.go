package repository

import (
	"context"
	"payment-service-provider/domain/entity"
)

type TransationRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, ID string) (*entity.Transaction, error)
}
