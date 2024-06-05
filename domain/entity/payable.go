package entity

import (
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

func NewPayable(ID string, transaction *Transaction) (*Payable, error) {
	p := &Payable{
		id:            ID,
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
	if transaction.paymentMethod.Method() == "credit_card" {
		p.status = waitFunds
		p.paymentDate = transaction.GetCreatedAt().AddDate(0, 0, 30)
	}
	return p, nil
}

func RestorePayable(id, clientID, transactionID string, status payableStatus, fee, value float32, createdAt, paymentDate time.Time) *Payable {
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
