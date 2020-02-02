package operations
import (
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

// TODO: fix this test
/*func TestDuplicatedOrder(t *testing.T){
	ordersPerIssuer := make(map[string][]order)
	ordersPerIssuerRef := &ordersPerIssuer
	ordersPerIssuerRefRef := &ordersPerIssuerRef
	order := order{}
	order.Timestamp = 1580060543
	order.Operation = "BUY"
	order.SharePrice = 10
	order.TotalShares = 10
	order.IssuerName = "ABC"
	if duplicatedOrder(order, ordersPerIssuerRefRef) {
		//t.Error("Order is not duplicated")
	}
	if !duplicatedOrder(order, &ordersPerIssuerRef) {
		//t.Error("Order is duplicated")
	}
	t.Fail()
}*/

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