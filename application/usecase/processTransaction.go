package usecase

import (
	"context"
	"fmt"
	"payment-service-provider/application/repository"
	"payment-service-provider/domain/entity"
	"time"

	"github.com/google/uuid"
)

type ProcessTransaction struct {
	transactionRepository repository.TransationRepository
	payableRepository     repository.PayableRepository
}

func NewProcessTransaction(
	transactionRepository repository.TransationRepository,
	payableRepository repository.PayableRepository,
) *ProcessTransaction {
	return &ProcessTransaction{
		transactionRepository: transactionRepository,
		payableRepository:     payableRepository,
	}
}

func (uc *ProcessTransaction) Execute(ctx context.Context, input *ProcessTransactionInput) (*ProcessTransactionOutput, error) {
	paymentMethod, ok := entity.PaymentMethods[input.PaymentMethod]
	if !ok {
		return nil, fmt.Errorf("invalid payment method %s", input.PaymentMethod)
	}
	card, err := entity.NewCard(input.Card.OwnerName, input.Card.VerificationCode, input.Card.Number, input.Card.ValidDate)
	if err != nil {
		return nil, err
	}
	transaction, err := entity.NewTransaction(uuid.New().String(), input.ClientID, input.Description, input.Value, time.Now(), card, paymentMethod)
	if err != nil {
		return nil, err
	}
	payable, err := entity.PayableFactory(uuid.New().String(), transaction)
	if err != nil {
		return nil, err
	}
	//TODO - UOW
	if err := uc.transactionRepository.Save(ctx, transaction); err != nil {
		return nil, err
	}
	if err := uc.payableRepository.Save(ctx, payable); err != nil {
		return nil, err
	}
	return NewProcessTransactionDTO(transaction, payable), nil
}

type ProcessTransactionInput struct {
	ClientID      string                      `json:"client_id"`
	Value         float32                     `json:"value"`
	Description   string                      `json:"description"`
	PaymentMethod string                      `json:"payment_method"`
	Card          ProcessTransactionCardInput `json:"card"`
}

type ProcessTransactionCardInput struct {
	Number           string    `json:"number"`
	VerificationCode string    `json:"verification_code"`
	OwnerName        string    `json:"owner_name"`
	ValidDate        time.Time `json:"valid_date"`
}

type ProcessTransactionOutput struct {
	TransactionDTO TransactionDTO `json:"transaction"`
	PayableDTO     PayableDTO     `json:"payable"`
}

type TransactionDTO struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"client_id"`
	Description   string    `json:"description"`
	Value         float32   `json:"value"`
	Card          CardDTO   `json:"card"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

type CardDTO struct {
	OwnerName        string    `json:"owner_name"`
	VerificationCode string    `json:"verification_code"`
	LastDigits       string    `json:"last_digits"`
	ValidDate        time.Time `json:"valid_date"`
}

type PayableDTO struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"client_id"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	FeeAmount     float32   `json:"fee_amount"`
	Amount        float32   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	PaymentDate   time.Time `json:"payment_date"`
}

func NewProcessTransactionDTO(transaction *entity.Transaction, payable *entity.Payable) *ProcessTransactionOutput {
	return &ProcessTransactionOutput{
		TransactionDTO: TransactionDTO{
			ID:          transaction.GetID(),
			ClientID:    transaction.GetClientID(),
			Description: transaction.GetDescription(),
			Value:       transaction.GetValue(),
			Card: CardDTO{
				OwnerName:        transaction.GetCard().GetOwnerName(),
				VerificationCode: transaction.GetCard().GetVerificationCode(),
				LastDigits:       transaction.GetCard().GetLastDigits(),
				ValidDate:        transaction.GetCard().GetValidDate(),
			},
		},
		PayableDTO: PayableDTO{
			ID:            payable.GetID(),
			ClientID:      payable.GetClientID(),
			TransactionID: payable.GetTransactionID(),
			Status:        payable.GetStatus(),
			FeeAmount:     payable.GetFeeAmount(),
			Amount:        payable.GetAmount(),
			CreatedAt:     payable.GetCreatedAt(),
			PaymentDate:   payable.GetPaymentDate(),
		},
	}
}
