package entity

type paymentMethod interface {
	String() string
}

type DebitCard struct {
	paymentMethod string
}

func (p *DebitCard) String() string {
	return p.paymentMethod
}

type CreditCard struct {
	paymentMethod string
}

func (p *CreditCard) String() string {
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
