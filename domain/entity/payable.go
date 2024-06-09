package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Payable interface {
	GetData() *PayableImpl
}

type PayableWithDebitCard struct {
	Data *PayableImpl
}

func (p *PayableWithDebitCard) GetData() *PayableImpl {
	return p.Data
}

func newPayableWithDebitCard(transaction *Transaction) (Payable, error) {
	p := &PayableImpl{
		id:            uuid.New().String(),
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		createdAt:     transaction.GetCreatedAt(),
	}
	feeCalculator, err := FeeCalculatorFactory(transaction.paymentMethod.Method())
	if err != nil {
		return nil, err
	}
	fee := feeCalculator.calculate(transaction.GetValue())
	p.feeAmount = fee
	p.amount = transaction.GetValue() - fee
	p.paymentDate = transaction.GetCreatedAt()
	p.status = paid
	return &PayableWithDebitCard{Data: p}, nil
}

type PayableWithCreditCard struct {
	Data *PayableImpl
}

func (p *PayableWithCreditCard) GetData() *PayableImpl {
	return p.Data
}

func newPayableWithCreditCard(transaction *Transaction) (Payable, error) {
	p := &PayableImpl{
		id:            uuid.New().String(),
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		createdAt:     transaction.GetCreatedAt(),
	}
	feeCalculator, err := FeeCalculatorFactory(transaction.paymentMethod.Method())
	if err != nil {
		return nil, err
	}
	fee := feeCalculator.calculate(transaction.GetValue())
	p.feeAmount = fee
	p.amount = transaction.GetValue() - fee
	p.status = waitFunds
	p.paymentDate = transaction.GetCreatedAt().AddDate(0, 0, 30)
	return &PayableWithCreditCard{Data: p}, nil
}

func PayableFactory(transaction *Transaction) (Payable, error) {
	if transaction.GetPaymentMethod() == "debit_card" {
		return newPayableWithDebitCard(transaction)
	}
	if transaction.GetPaymentMethod() == "credit_card" {
		return newPayableWithCreditCard(transaction)
	}
	return nil, fmt.Errorf("invalid transaciton paymentMethod: %s", transaction.GetPaymentMethod())
}

type PayableImpl struct {
	id            string
	clientID      string
	transactionID string
	status        payableStatus
	feeAmount     float32
	amount        float32
	createdAt     time.Time
	paymentDate   time.Time
}

func (p *PayableImpl) GetID() string {
	return p.id
}

func (p *PayableImpl) GetClientID() string {
	return p.clientID
}

func (p *PayableImpl) GetTransactionID() string {
	return p.transactionID
}

func (p *PayableImpl) GetStatus() string {
	return p.status.String()
}

func (p *PayableImpl) GetFeeAmount() float32 {
	return p.feeAmount
}

func (p *PayableImpl) GetAmount() float32 {
	return p.amount
}

func (p *PayableImpl) GetCreatedAt() time.Time {
	return p.createdAt
}

func (p *PayableImpl) GetPaymentDate() time.Time {
	return p.paymentDate
}

func RecoverPayableData(id, clientID, transactionID string, status payableStatus, fee, value float32, createdAt, paymentDate time.Time) *PayableImpl {
	p := &PayableImpl{
		id:            id,
		clientID:      clientID,
		transactionID: transactionID,
		status:        status,
		feeAmount:     fee,
		amount:        value,
		createdAt:     createdAt,
		paymentDate:   paymentDate,
	}
	return p
}

func RestorePayable(p *PayableImpl) Payable {
	if p.status == paid {
		return &PayableWithDebitCard{Data: p}
	}
	return &PayableWithCreditCard{Data: p}
}

type payableStatus string

const (
	paid      payableStatus = "paid"
	waitFunds payableStatus = "wait_funds"
)

func (p payableStatus) String() string {
	return string(p)
}

func GetPayableStatus(s string) payableStatus {
	return payableStatuses[s]
}

var payableStatuses = map[string]payableStatus{
	"paid":       paid,
	"wait_funds": waitFunds,
}
