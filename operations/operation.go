package operations

import (
	"strings"
	"sync"
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

// Operation - holds the orders from a json file
type Operation struct {
	InitialBalance balance `json:"initialBalances"`
	Orders         []order `json:"orders"`
}

type businessError struct {
	ErrorType   string
	OrderFailed order
}

// Output - Holds the result from running json file operations
type Output struct {
	CurrentBalance balance         `json:"currentBalance"`
	BusinessErrors []businessError `json:"businessErrors"`
}

// isOrderInvalid - verify if the structure of the order is valid, type and not null values
func isOrderInvalid(order order) bool {
	return order.Timestamp < 0 || order.TotalShares < 0 || order.SharePrice < 0 || len(order.IssuerName) <= 0 || (order.Operation != "BUY" && order.Operation != "SELL")
}

// validMarketHoursOperation - check if the order is between 06:00 and 15:00 hours
func validMarketHoursOperation(timestamp int64) bool {
	date := time.Unix(timestamp, 0)
	totalSeconds := date.Second() + date.Minute()*60 + date.Hour()*3600
	return totalSeconds > 21600 && totalSeconds < 54000
}

/* duplicatedOrder - validates if an orders is duplicated. the criteria is:
	- same issuer, shares and operation
	- the operation is in 5 minutes difference between other operations
*/
func duplicatedOrder(currentOrder order, ordersPerIssuer **map[string][]order) bool {
	orders, exists := (**ordersPerIssuer)[currentOrder.IssuerName]
	if !exists {
		(**ordersPerIssuer)[currentOrder.IssuerName] = []order{currentOrder}
		return false
	}

	for _, order := range orders {
		if currentOrder.TotalShares == order.TotalShares &&
			currentOrder.SharePrice == order.SharePrice &&
			currentOrder.Operation == order.Operation {
			if currentOrder.Timestamp == order.Timestamp {
				return true
			}
			if order.Timestamp > currentOrder.Timestamp {
				return order.Timestamp-currentOrder.Timestamp <= 300
			}
			if currentOrder.Timestamp > order.Timestamp {
				return currentOrder.Timestamp-order.Timestamp <= 300
			}
		}
	}
	(**ordersPerIssuer)[currentOrder.IssuerName] = append((**ordersPerIssuer)[currentOrder.IssuerName], currentOrder)
	return false
}

// canSell - validates the shares of the issuer to check amount
func canSell(order order, balance balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return issuer.TotalShares >= order.TotalShares, operationCost, -order.TotalShares
		}
	}
	return false, 0, 0
}

// canBuy - validates the cash of the issuer to check amount
func canBuy(order order, balance balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return balance.Cash >= operationCost, -operationCost, order.TotalShares
		}
	}
	return false, 0, 0
}

// updateBalance - depending if its buy/sell update share and cash accordingly
func updateBalance(operation *Operation, operationIssuer string, cost int, shares int) {
	operation.InitialBalance.Cash += cost
	for i, issuer := range operation.InitialBalance.Issuers {
		if issuer.IssuerName == operationIssuer {
			operation.InitialBalance.Issuers[i].TotalShares += shares
			break
		}
	}
}

// runOrder - store the order in a map of issuers, for duplicated validations, update balance according of the order type
func runOrder(operation *Operation, order order, ordersPerIssuer *map[string][]order, output *Output) {
	if isOrderInvalid(order) {
		bError := businessError{}
		bError.ErrorType = "INVALID OPERATION"
		bError.OrderFailed = order
		output.BusinessErrors = append(output.BusinessErrors, bError)
		return
	}
	if !validMarketHoursOperation(order.Timestamp) {
		bError := businessError{}
		bError.ErrorType = "CLOSED MARKET"
		bError.OrderFailed = order
		output.BusinessErrors = append(output.BusinessErrors, bError)
		return
	}
	if duplicatedOrder(order, &ordersPerIssuer) {
		bError := businessError{}
		bError.ErrorType = "DUPLICATED OPERATION"
		bError.OrderFailed = order
		output.BusinessErrors = append(output.BusinessErrors, bError)
		return
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

// PerformOperation - run all orders in a json file. Uses go routines per file for better performance
func PerformOperation(operation *Operation, wg *sync.WaitGroup) Output {
	defer wg.Done()
	output := Output{}
	ordersPerIssuer := make(map[string][]order)
	for _, order := range operation.Orders {
		runOrder(operation, order, &ordersPerIssuer, &output)
	}
	output.CurrentBalance.Cash = operation.InitialBalance.Cash
	output.CurrentBalance.Issuers = operation.InitialBalance.Issuers
	return output
}