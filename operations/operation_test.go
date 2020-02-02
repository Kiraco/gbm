package operations

import (
	"sync"
	"testing"
)

func TestIsOrderInvalidFalse(t *testing.T){
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if isOrderInvalid(order) {
		t.Error("Order should be valid")
	}
}

func TestIsOrderInvalidTimestamp(t *testing.T){
	order := order{}
	order.Timestamp = -1
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if !isOrderInvalid(order) {
		t.Error("Order should be invalid")
	}
}

func TestIsOrderInvalidOperation(t *testing.T){
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "B"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if !isOrderInvalid(order) {
		t.Error("Order should be invalid")
	}
}

func TestIsOrderInvalidSharePrice(t *testing.T){
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = -1
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if !isOrderInvalid(order) {
		t.Error("Order should be invalid")
	}
}

func TestIsOrderInvalidTotalShares(t *testing.T){
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = -1
	order.IssuerName = "ABC"
	if !isOrderInvalid(order) {
		t.Error("Order should be invalid")
	}
}

func TestIsOrderInvalidIssuerName(t *testing.T){
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	if !isOrderInvalid(order) {
		t.Error("Order should be invalid")
	}
}

func TestValidMarketHoursOperation(t *testing.T){
	timestamp := 1580060543
	if !validMarketHoursOperation(int64(timestamp)) {
		t.Error("timestamp should be valid")
	}
}

func TestValidMarketHoursOperationFalseEarly(t *testing.T){
	timestamp := 1580558399
	if validMarketHoursOperation(int64(timestamp)) {
		t.Error("timestamp should be valid")
	}
}

func TestValidMarketHoursOperationFalseLate(t *testing.T){
	timestamp := 1580590801
	if validMarketHoursOperation(int64(timestamp)) {
		t.Error("timestamp should be invalid")
	}
}

func TestDuplicatedOrderFalse(t *testing.T){
	ordersPerIssuer := make(map[string][]order)
	ordersPerIssuerRef := &ordersPerIssuer
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is not duplicated")
	}
}

func TestDuplicatedOrderSameHour(t *testing.T){
	ordersPerIssuer := make(map[string][]order)
	ordersPerIssuerRef := &ordersPerIssuer
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is not duplicated")
	}
	if !duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is duplicated")
	}
}

func TestDuplicatedOrderLess5Minutes(t *testing.T){
	ordersPerIssuer := make(map[string][]order)
	ordersPerIssuerRef := &ordersPerIssuer
	order := order{}
	order.Timestamp = 1580060842
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is not duplicated")
	}
	order.Timestamp = 1580060543
	if !duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is duplicated")
	}
}

func TestDuplicatedOrder5Minutes(t *testing.T){
	ordersPerIssuer := make(map[string][]order)
	ordersPerIssuerRef := &ordersPerIssuer
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is not duplicated")
	}
	order.Timestamp = 1580060843
	if !duplicatedOrder(order, &ordersPerIssuerRef) {
		t.Error("Order is duplicated")
	}
}

func TestCanSell(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	can, _, _ := canSell(order, balance)
	if !can {
		t.Error("Order can be sell")
	}
}

func TestCanSellFalse(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 200
	order.IssuerName = "ABC"

	can, _, _ := canSell(order, balance)
	if can {
		t.Error("Order can't be sell")
	}
}

func TestCanSellFalseNoIssuer(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "NON_EXISTENT_ISSUER_NAME"

	can, _, _ := canSell(order, balance)
	if can {
		t.Error("Order can't be sell")
	}
}

func TestCanBuy(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	can, _, _ := canBuy(order, balance)
	if !can {
		t.Error("Order can be sell")
	}
}

func TestCanBuyFalse(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 200
	order.IssuerName = "ABC"

	can, _, _ := canSell(order, balance)
	if can {
		t.Error("Order can't be sell")
	}
}

func TestCanBuyFalseNoIssuer(t *testing.T){
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = 1000
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "NON_EXISTENT_ISSUER_NAME"

	can, _, _ := canSell(order, balance)
	if can {
		t.Error("Order can't be sell")
	}
}

func TestUpdateBalance(t *testing.T) {
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	updateBalance(&operation, order.IssuerName, order.TotalShares * order.SharePrice, order.TotalShares)
	if operation.InitialBalance.Cash == initialCash {
		t.Error("Cash should be updated")
	}
}

func TestRunOrderBuy(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if output.CurrentBalance.Cash == initialCash {
		t.Error("Cash should be updated, valid buy order")
	}
}

func TestRunOrderSell(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if output.CurrentBalance.Cash == initialCash {
		t.Error("Cash should be updated, valid buy order")
	}
}

func TestRunOrderCantBuy(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 1000
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if len(output.BusinessErrors) == 0 {
		t.Error("Should be an error since can't buy")
	}
}

func TestRunOrderCantSell(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 1000
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if len(output.BusinessErrors) == 0 {
		t.Error("Should be an error since can't sell")
	}
}

func TestRunOrderInvalidOrder(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = ""
	order.SharePrice = 10
	order.TotalShares = 1000
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if len(output.BusinessErrors) == 0 {
		t.Error("Should be an error since its and invalid operation")
	}
}

func TestRunOrderClosedMarket(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580558399
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 1000
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	if len(output.BusinessErrors) == 0 {
		t.Error("Should be an error since timestamp is in invalid hour")
	}
}

func TestRunOrderDuplicatedOrder(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	ordersPerIssuer := make(map[string][]order)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)
	output := Output{}

	runOrder(&operation, order, &ordersPerIssuer, &output)
	runOrder(&operation, order, &ordersPerIssuer, &output)
	if len(output.BusinessErrors) == 0 {
		t.Error("Should be an error since there is a duplicated order")
	}
}

func TestPerformOperation(t *testing.T){
	initialCash := 1000
	issuers := []issuer{
		{IssuerName: "ABC", TotalShares: 100, SharePrice: 10},
	}
	balance := balance{}
	balance.Cash = initialCash
	balance.Issuers = issuers

	wg := sync.WaitGroup{}
	wg.Add(1)

	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "SELL"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"

	operation := Operation{}
	operation.InitialBalance = balance
	operation.Orders = append(operation.Orders, order)

	output := PerformOperation(&operation, &wg)
	wg.Wait()
	if output.CurrentBalance.Cash == initialCash {
		t.Error("Should be an error since orders from operation weren't run")
	}
}