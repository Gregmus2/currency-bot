package main

type Receiver struct {
	buy  float64
	sell float64
}

func (r *Receiver) Buy() float64 {
	return r.buy
}

func (r *Receiver) Sell() float64 {
	return r.sell
}

func (r *Receiver) BankName() string {
	return ""
}
