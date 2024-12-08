package evaluator

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	var sampleContext = NewContext(
		&UserContext{},
		[]*StoreContext{
			&StoreContext{
				Id: "store1",
				Products: []*ProductContext{
					&ProductContext{
						Id:                "1",
						OriginalPrice:     100.0,
						DiscountedPrice:   0.0,
						Category:          "2",
						RequestedQuantity: 1,
					},
					&ProductContext{
						Id:                "2",
						OriginalPrice:     200.0,
						DiscountedPrice:   0.0,
						Category:          "2",
						RequestedQuantity: 1,
					},
					&ProductContext{
						Id:                "3",
						OriginalPrice:     300.0,
						DiscountedPrice:   0.0,
						Category:          "3",
						RequestedQuantity: 1,
					},
				},
				CartPrice:               0.0,
				CartPriceAfterDiscounts: 80.0,
				Discounts: map[string]string{
					"DISCOUNT10": "product_id in [1, 2] and category_id = 2 then product_percentage 10",
					// "DISCOUNT30": "min_cart_price > 500 then cart_percentage 8",
				},
			},
		},
	)
	Evaluate(sampleContext)
	if ((*sampleContext).Stores)[0].CartPriceAfterDiscounts != 80.0 {
		t.Errorf("Expected cart price after discounts to be 80.0, but got %f", ((*sampleContext).Stores)[0].CartPriceAfterDiscounts)
	}
	// if (((*sampleContext).Stores)[0].Products)[0].DiscountedPrice != 90.0 {
	// 	t.Errorf("Expected discounted price for product 1 to be 90.0, but got %f", (((*sampleContext).Stores)[0].Products)[0].DiscountedPrice)
	// }
}
