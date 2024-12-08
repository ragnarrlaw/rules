package evaluator

type Context struct {
	User   *UserContext
	Stores []*StoreContext
}

type UserContext struct {
}

type StoreContext struct {
	Id                      string
	Products                []*ProductContext
	CartPrice               float64
	CartPriceAfterDiscounts float64
	Discounts               map[string]string
}

type ProductContext struct {
	Id                string
	Category          string
	RequestedQuantity float64
	OriginalPrice     float64
	DiscountedPrice   float64
}

func NewContext(userData *UserContext, storeData []*StoreContext) *Context {
	return &Context{userData, storeData}
}
