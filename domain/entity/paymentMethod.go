package entity

type paymentMethod interface {
	Method() string
}

type DebitCard struct {
	paymentMethod string
}

func (p *DebitCard) Method() string {
	return p.paymentMethod
}

type CreditCard struct {
	paymentMethod string
}

func (p *CreditCard) Method() string {
	return p.paymentMethod
}

var debitCard *DebitCard = &DebitCard{
	paymentMethod: "debit_card",
}

var creditCard *CreditCard = &CreditCard{
	paymentMethod: "credit_card",
}

var PaymentMethods = map[string]paymentMethod{
	"debit_card":  debitCard,
	"credit_card": creditCard,
}
