package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	//"os"
)

func itHasEnoughCash(cash int, totalShares int, sharePrice int) bool {
	return (totalShares * sharePrice) < cash

}

func canSell(order Order, balance Balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return issuer.TotalShares >= order.TotalShares, operationCost, -order.TotalShares
		}
	}
	return false, 0, 0
}

func canBuy(order Order, balance Balance) (can bool, cost int, shares int) {
	for _, issuer := range balance.Issuers {
		if issuer.IssuerName == order.IssuerName {
			operationCost := order.SharePrice * order.TotalShares
			return balance.Cash >= operationCost, -operationCost, order.TotalShares
		}
	}
	return false, 0, 0
}

func validTimeOperation(timestamp int64) bool {
	date := time.Unix(timestamp, 0)
	totalSeconds := date.Second() + date.Minute() * 60 + date.Hour() * 3600
	return totalSeconds > 21600 && totalSeconds < 54000
}

func loadData(filePath string) Operation {
	file, _ := ioutil.ReadFile(filePath)
	operation := Operation{}
	err := json.Unmarshal(file, &operation)
	if err != nil {
		panic(err)
	}
	return operation
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

func performOperation(operation *Operation) Output {
	output := Output{}

	for _, order := range operation.Orders {
		if !validTimeOperation(order.Timestamp) {
			error := BusinessError{}
			error.ErrorType = "CLOSED MARKET"
			error.OrderFailed = order
			output.BusinessErrors = append(output.BusinessErrors, error)
			continue
		}
		switch strings.ToUpper(order.Operation){
			case "BUY":
				can, cost, shares := canBuy(order, operation.InitialBalance )
				if can {
					updateBalance(operation, order.IssuerName, cost, shares)
				} else {
					error := BusinessError{}
					error.ErrorType = "INSUFFICIENT BALANCE"
					error.OrderFailed = order
					output.BusinessErrors = append(output.BusinessErrors, error)
				}
				break
		case "SELL":
			can, cost, shares := canSell(order, operation.InitialBalance )
			if can {
				updateBalance(operation, order.IssuerName, cost, shares)
			} else {
				error := BusinessError{}
				error.ErrorType = "INSUFFICIENT STOCKS"
				error.OrderFailed = order
				output.BusinessErrors = append(output.BusinessErrors, error)
			}
			break
		}
	}

	output.CurrentBalance.Cash = operation.InitialBalance.Cash
	output.CurrentBalance.Issuers = operation.InitialBalance.Issuers
	return output
}

func main() {
	path := "/Users/donovan/Documents/Personal/Projects/go/gbm/test.json"
	//filepath := os.Args[1]
	operation := loadData(path)
	fmt.Printf("%+v\n", operation)
	output := performOperation(&operation)
	fmt.Printf("%+v\n", output)
}
