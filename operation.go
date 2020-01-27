package main

type Order struct {
	Timestamp   int64    `json:"timestamp`
	Operation   string `json:"operation"`
	IssuerName  string `json:"IssuerName"`
	TotalShares int    `json:"TotalShares`
	SharePrice  int    `json:"SharePrice`
}

type Issuer struct {
	IssuerName  string `json:"issuerName"`
	TotalShares int    `json:"totalShares"`
	SharePrice  int    `json:"sharePrice"`
}

type Balance struct {
	Cash    int      `json:"cash"`
	Issuers []Issuer `json:"issuers`
}

type Operation struct {
	InitialBalance Balance `json:"initialBalances"`
	Orders         []Order `json:"orders"`
}

type BusinessError struct {
	ErrorType	string
	OrderFailed	Order
}

type Output struct {
	CurrentBalance Balance `json:"currentBalance"`
	BusinessErrors []BusinessError `json:""businessErrors`

}
