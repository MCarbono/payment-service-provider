package entity

import (
	"fmt"
	"time"
)

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

func newPayableWithDebitCard(ID string, transaction *Transaction) *Payable {
	fee := calculateFee(transaction.value, transaction.paymentMethod.Fee())
	return &Payable{
		id:            ID,
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		status:        paid,
		feeAmount:     fee,
		amount:        transaction.GetValue() - fee,
		createdAt:     transaction.GetCreatedAt(),
		paymentDate:   transaction.GetCreatedAt(),
	}
}

func newPayableWithCreditCard(ID string, transaction *Transaction) *Payable {
	fee := calculateFee(transaction.value, transaction.paymentMethod.Fee())
	return &Payable{
		id:            ID,
		clientID:      transaction.clientID,
		transactionID: transaction.GetID(),
		status:        waitFunds,
		feeAmount:     fee,
		amount:        transaction.GetValue() - fee,
		createdAt:     transaction.GetCreatedAt(),
		paymentDate:   transaction.GetCreatedAt().AddDate(0, 0, 30),
	}
}

func PayableFactory(ID string, transaction *Transaction) (*Payable, error) {
	if transaction.paymentMethod.Method() == "debit_card" {
		return newPayableWithDebitCard(ID, transaction), nil
	}
	if transaction.paymentMethod.Method() == "credit_card" {
		return newPayableWithCreditCard(ID, transaction), nil
	}
	return nil, fmt.Errorf("invalid paymentMethod: %s", transaction.paymentMethod.Method())
}

func NewPayable(id, clientID, transactionID string, status payableStatus, fee, value float32, createdAt, paymentDate time.Time) *Payable {
	return &Payable{
		id:            id,
		clientID:      clientID,
		transactionID: transactionID,
		status:        status,
		feeAmount:     fee,
		amount:        value,
		createdAt:     createdAt,
		paymentDate:   paymentDate,
	}
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
