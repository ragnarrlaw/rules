package evaluator

import "testing"

func TestStoreContextDiscounts(t *testing.T) {
	discounts := map[string]string{
		"DISCOUNT10": "10 off",
		"DISCOUNT20": "20 off",
	}

	storeContext := &StoreContext{
		Id:                      "store1",
		Products:                []*ProductContext{},
		CartPrice:               100.0,
		CartPriceAfterDiscounts: 80.0,
		Discounts:               discounts,
	}

	if storeContext.Discounts == nil {
		t.Errorf("Expected Discounts to be initialized, but got nil")
	}

	if len(storeContext.Discounts) != 2 {
		t.Errorf("Expected 2 discounts, but got %d", len(storeContext.Discounts))
	}

	if (storeContext.Discounts)["DISCOUNT10"] != "10 off" {
		t.Errorf("Expected '10 off' for DISCOUNT10, but got %s", (storeContext.Discounts)["DISCOUNT10"])
	}

	if (storeContext.Discounts)["DISCOUNT20"] != "20 off" {
		t.Errorf("Expected '20 off' for DISCOUNT20, but got %s", (storeContext.Discounts)["DISCOUNT20"])
	}
}
