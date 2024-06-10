package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BuildProcess interface {
	SetID() BuildProcess
	SetClientId(clientID string) BuildProcess
	SetTransactionID(transactionID string) BuildProcess
	SetStatus() BuildProcess
	SetFeeAmount(amount float32) BuildProcess
	SetAmount(amount float32) BuildProcess
	SetCreatedAt(createdAt time.Time) BuildProcess
	SetPaymentDate(paymentDate time.Time) BuildProcess
	Build() Payable
}

type PayableBuilderWithDebitCard struct {
	p             Payable
	feePercentage float32
}

func NewPayableBuilderWithDebitCard() *PayableBuilderWithDebitCard {
	return &PayableBuilderWithDebitCard{
		feePercentage: 0.03,
	}
}

func (b *PayableBuilderWithDebitCard) SetID() BuildProcess {
	b.p.id = uuid.NewString()
	return b
}

func (b *PayableBuilderWithDebitCard) SetClientId(clientID string) BuildProcess {
	b.p.clientID = clientID
	return b
}

func (b *PayableBuilderWithDebitCard) SetTransactionID(transactionID string) BuildProcess {
	b.p.transactionID = transactionID
	return b
}

func (b *PayableBuilderWithDebitCard) SetStatus() BuildProcess {
	b.p.status = paid
	return b
}

func (b *PayableBuilderWithDebitCard) SetFeeAmount(amount float32) BuildProcess {
	b.p.feeAmount = amount * b.feePercentage
	return b
}

func (b *PayableBuilderWithDebitCard) SetAmount(amount float32) BuildProcess {
	b.p.amount = amount - b.p.feeAmount
	return b
}

func (b *PayableBuilderWithDebitCard) SetCreatedAt(createdAt time.Time) BuildProcess {
	b.p.createdAt = createdAt
	return b
}

func (b *PayableBuilderWithDebitCard) SetPaymentDate(paymentDate time.Time) BuildProcess {
	b.p.paymentDate = paymentDate
	return b
}

func (b *PayableBuilderWithDebitCard) Build() Payable {
	return b.p
}

type PayableBuilderWithCreditCard struct {
	p             Payable
	feePercentage float32
}

func NewPayableBuilderWithCreditCard() *PayableBuilderWithCreditCard {
	return &PayableBuilderWithCreditCard{
		feePercentage: 0.05,
	}
}

func (b *PayableBuilderWithCreditCard) SetID() BuildProcess {
	b.p.id = uuid.NewString()
	return b
}

func (b *PayableBuilderWithCreditCard) SetClientId(clientID string) BuildProcess {
	b.p.clientID = clientID
	return b
}

func (b *PayableBuilderWithCreditCard) SetTransactionID(transactionID string) BuildProcess {
	b.p.transactionID = transactionID
	return b
}

func (b *PayableBuilderWithCreditCard) SetStatus() BuildProcess {
	b.p.status = waitFunds
	return b
}

func (b *PayableBuilderWithCreditCard) SetFeeAmount(amount float32) BuildProcess {
	b.p.feeAmount = amount * b.feePercentage
	return b
}

func (b *PayableBuilderWithCreditCard) SetAmount(amount float32) BuildProcess {
	b.p.amount = amount - b.p.feeAmount
	return b
}

func (b *PayableBuilderWithCreditCard) SetCreatedAt(createdAt time.Time) BuildProcess {
	b.p.createdAt = createdAt
	return b
}

func (b *PayableBuilderWithCreditCard) SetPaymentDate(paymentDate time.Time) BuildProcess {
	b.p.paymentDate = paymentDate.AddDate(0, 0, 30)
	return b
}

func (b *PayableBuilderWithCreditCard) Build() Payable {
	return b.p
}

type payableDirector struct {
	builder BuildProcess
}

func (f *payableDirector) Construct(transaction *Transaction) {
	f.builder.SetID().SetTransactionID(transaction.id).
		SetClientId(transaction.clientID).SetStatus().
		SetFeeAmount(transaction.value).SetAmount(transaction.value).
		SetCreatedAt(transaction.createdAt).SetPaymentDate(transaction.createdAt)
}

func (f *payableDirector) SetBuilder(b BuildProcess) {
	f.builder = b
}

var payableDirectorSingleton *payableDirector

func GetPayableDirector() *payableDirector {
	if payableDirectorSingleton == nil {
		payableDirectorSingleton = &payableDirector{}
		return payableDirectorSingleton
	} else {
		return payableDirectorSingleton
	}
}

func BuilderFactory(transaction *Transaction) (BuildProcess, error) {
	if transaction.paymentMethod == debitCard {
		return NewPayableBuilderWithDebitCard(), nil
	}
	if transaction.paymentMethod == creditCard {
		return NewPayableBuilderWithCreditCard(), nil
	}
	return nil, fmt.Errorf("invalid transaction paymentMethod: %s", transaction.PaymentMethod())
}
