package usecase

// import (
// 	"context"
// 	"fmt"
// 	"payment-service-provider/application/uow"
// 	"payment-service-provider/domain/entity"
// 	"time"

// 	"github.com/google/uuid"
// )

// type ProcessTransactionBuilder struct {
// 	uow uow.UnitOfWork
// }

// func NewProcessTransactionBuilder(
// 	uow uow.UnitOfWork,
// ) *ProcessTransaction {
// 	return &ProcessTransaction{
// 		uow: uow,
// 	}
// }

// func (uc *ProcessTransactionBuilder) Execute(ctx context.Context, input *ProcessTransactionInput) (*ProcessTransactionOutput, error) {
// 	paymentMethod, ok := entity.PaymentMethods[input.PaymentMethod]
// 	if !ok {
// 		return nil, fmt.Errorf("invalid payment method %s", input.PaymentMethod)
// 	}
// 	card, err := entity.NewCard(input.Card.OwnerName, input.Card.VerificationCode, input.Card.Number, input.Card.ValidDate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	transaction, err := entity.NewTransaction(uuid.New().String(), input.ClientID, input.Description, input.Value, time.Now(), card, paymentMethod)
// 	if err != nil {
// 		return nil, err
// 	}
// 	director := entity.GetPayableDirector()
// 	builder, err := entity.BuilderFactory(transaction)
// 	if err != nil {
// 		return nil, err
// 	}
// 	director.SetBuilder(builder)
// 	director.Construct(transaction)
// 	payable := builder.Build()
// 	err = uc.uow.RunTx(ctx, func(ctx context.Context, repositories uow.Repositories) error {
// 		err := repositories.Transaction().Save(ctx, transaction)
// 		if err != nil {
// 			return err
// 		}
// 		err = repositories.Payable().Save(ctx, payable)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return NewProcessTransactionDTO(transaction, payable), err
// }
