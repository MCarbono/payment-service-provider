package entity

type paymentMethod interface {
	Method() string
	Fee() float32
}

type DebitCard struct {
	fee           float32
	paymentMethod string
}

func (p *DebitCard) Method() string {
	return p.paymentMethod
}

func (p *DebitCard) Fee() float32 {
	return p.fee
}

type CreditCard struct {
	fee           float32
	paymentMethod string
}

func (p *CreditCard) Method() string {
	return p.paymentMethod
}

func (p *CreditCard) Fee() float32 {
	return p.fee
}

var debitCard *DebitCard = &DebitCard{
	fee:           0.03,
	paymentMethod: "debit_card",
}

var creditCard *CreditCard = &CreditCard{
	fee:           0.05,
	paymentMethod: "credit_card",
}

var PaymentMethods = map[string]paymentMethod{
	"debit_card":  debitCard,
	"credit_card": creditCard,
}
