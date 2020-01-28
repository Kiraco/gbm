package operations

import (
	"strings"
	"time"
)

type order struct {
	Timestamp   int64  `json:"timestamp"`
	Operation   string `json:"operation"`
	IssuerName  string `json:"IssuerName"`
	TotalShares int    `json:"TotalShares"`
	SharePrice  int    `json:"SharePrice"`
}

type issuer struct {
	IssuerName  string `json:"issuerName"`
	TotalShares int    `json:"totalShares"`
	SharePrice  int    `json:"sharePrice"`
}

type balance struct {
	Cash    int      `json:"cash"`
	Issuers []issuer `json:"issuers"`
}

type Operation struct {
	InitialBalance balance `json:"initialBalances"`
	Orders         []order `json:"orders"`
}

type businessError struct {
	ErrorType   string
	OrderFailed order
}

type Output struct {
	CurrentBalance balance         `json:"currentBalance"`
	BusinessErrors []businessError `json:"businessErrors"`
}

func canSell(order order, balance balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return issuer.TotalShares >= order.TotalShares, operationCost, -order.TotalShares
		}
	}
	return false, 0, 0
}

func canBuy(order order, balance balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return balance.Cash >= operationCost, -operationCost, order.TotalShares
		}
	}
	return false, 0, 0
}

func validMarketHoursOperation(timestamp int64) bool {
	date := time.Unix(timestamp, 0)
	totalSeconds := date.Second() + date.Minute()*60 + date.Hour()*3600
	return totalSeconds > 21600 && totalSeconds < 54000
}

func updateBalance(operation *Operation, operationIssuer string, cost int, shares int) {
	operation.InitialBalance.Cash += cost
	for i, issuer := range operation.InitialBalance.Issuers {
		if issuer.IssuerName == operationIssuer {
			operation.InitialBalance.Issuers[i].TotalShares += shares
			break
		}
	}
}

func duplicatedOperation(previousTime int64, currentTime int64) bool {
	previous := time.Unix(previousTime, 0)
	current := time.Unix(currentTime, 0)
	previous.Add(time.Minute * 5)
	return previous.Equal(current) || previous.After(current)
}

func PerformOperation(operation *Operation) Output {
	output := Output{}

	for i, order := range operation.Orders {
		if !validMarketHoursOperation(order.Timestamp) {
			bError := businessError{}
			bError.ErrorType = "CLOSED MARKET"
			bError.OrderFailed = order
			output.BusinessErrors = append(output.BusinessErrors, bError)
			continue
		}
		if i != 0 && duplicatedOperation(operation.Orders[i-1].Timestamp, order.Timestamp) {
			bError := businessError{}
			bError.ErrorType = "DUPLICATED OPERATION"
			bError.OrderFailed = order
			output.BusinessErrors = append(output.BusinessErrors, bError)
			continue
		}
		switch strings.ToUpper(order.Operation) {
		case "BUY":
			can, cost, shares := canBuy(order, operation.InitialBalance)
			if can {
				updateBalance(operation, order.IssuerName, cost, shares)
			} else {
				bError := businessError{}
				bError.ErrorType = "INSUFFICIENT BALANCE"
				bError.OrderFailed = order
				output.BusinessErrors = append(output.BusinessErrors, bError)
			}
			break
		case "SELL":
			can, cost, shares := canSell(order, operation.InitialBalance)
			if can {
				updateBalance(operation, order.IssuerName, cost, shares)
			} else {
				bError := businessError{}
				bError.ErrorType = "INSUFFICIENT STOCKS"
				bError.OrderFailed = order
				output.BusinessErrors = append(output.BusinessErrors, bError)
			}
			break
		}
	}

	output.CurrentBalance.Cash = operation.InitialBalance.Cash
	output.CurrentBalance.Issuers = operation.InitialBalance.Issuers
	return output
}
