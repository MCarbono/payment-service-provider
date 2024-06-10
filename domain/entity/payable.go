package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PayableInterface interface {
	GetData() *Payable
}

type PayableWithDebitCard struct {
	Data *Payable
}

func (p *PayableWithDebitCard) GetData() *Payable {
	return p.Data
}

func newPayableWithDebitCard(transaction *Transaction) (PayableInterface, error) {
	p := &Payable{
		id:            uuid.New().String(),
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		createdAt:     transaction.GetCreatedAt(),
	}
	feeCalculator, err := FeeCalculatorFactory(transaction.paymentMethod)
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
	Data *Payable
}

func (p *PayableWithCreditCard) GetData() *Payable {
	return p.Data
}

func newPayableWithCreditCard(transaction *Transaction) (PayableInterface, error) {
	p := &Payable{
		id:            uuid.New().String(),
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		createdAt:     transaction.GetCreatedAt(),
	}
	feeCalculator, err := FeeCalculatorFactory(transaction.paymentMethod)
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

func PayableFactory(transaction *Transaction) (PayableInterface, error) {
	if transaction.PaymentMethod() == debitCard {
		return newPayableWithDebitCard(transaction)
	}
	if transaction.PaymentMethod() == creditCard {
		return newPayableWithCreditCard(transaction)
	}
	return nil, fmt.Errorf("invalid transaciton paymentMethod: %s", transaction.PaymentMethod())
}

type Payable struct {
	id            string
	clientID      string
	transactionID string
	status        payableStatus
	feeAmount     float32
	amount        float32
	createdAt     time.Time
	paymentDate   time.Time
}

func (p *Payable) GetID() string {
	return p.id
}

func (p *Payable) GetClientID() string {
	return p.clientID
}

func (p *Payable) GetTransactionID() string {
	return p.transactionID
}

func (p *Payable) GetStatus() string {
	return p.status.String()
}

func (p *Payable) GetFeeAmount() float32 {
	return p.feeAmount
}

func (p *Payable) GetAmount() float32 {
	return p.amount
}

func (p *Payable) GetCreatedAt() time.Time {
	return p.createdAt
}

func (p *Payable) GetPaymentDate() time.Time {
	return p.paymentDate
}

func RecoverPayableData(id, clientID, transactionID string, status payableStatus, fee, value float32, createdAt, paymentDate time.Time) *Payable {
	p := &Payable{
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

func RestorePayable(p *Payable) PayableInterface {
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
