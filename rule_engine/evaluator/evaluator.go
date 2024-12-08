package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ragnarrlaw/rules/rule_engine/lexer"
	"github.com/ragnarrlaw/rules/rule_engine/parser"
)

func Evaluate(c *Context) error {
	for _, store := range c.Stores {
		// add a discount type -> product, category, and cart(store wide)
		for _, dRule := range store.Discounts {
			l := lexer.NewLexer(dRule)
			p := parser.NewParser(l)
			r, err := p.ParseRule()
			if err != nil {
				return err
			}
			for _, product := range store.Products {
				result, err := EvaluateLogicalCondition(r.Condition, store, product)
				if err != nil {
					return err
				}
				if result {
					ApplyAction(r.Action, store, product)
				}
				store.CartPrice += product.OriginalPrice * product.RequestedQuantity
			}
		}
	}
	return nil
}

func EvaluateLogicalCondition(r *parser.LogicalCondition, store *StoreContext, product *ProductContext) (bool, error) {
	leftResult := evaluateSingleCondition(r.Left, store, product)
	rightResult := evaluateSingleCondition(r.Right, store, product)
	switch r.Operator {
	case "AND":
		return leftResult && rightResult, nil
	case "OR":
		return leftResult || rightResult, nil
	default:
		return false, fmt.Errorf("unknown logical operator %s", r.Operator)
	}
}

// TODO: ADD EVALUATION FUNCTIONS TO EACH DISCOUNT TYPE -> PRODUCT, CATEGORY, AND CART(STORE-WIDE)
func evaluateSingleCondition(c *parser.Condition, store *StoreContext, product *ProductContext) bool {
	switch strings.ToLower(c.Key) {
	case "product_id":
		return compare(c.Operator, c.Value, product.Id)
	case "category_id":
		return compare(c.Operator, c.Value, product.Category)
	case "requested_quantity":
		return compare(c.Operator, c.Value, product.RequestedQuantity)
	case "cart_price":
		return compare(c.Operator, c.Value, store.CartPrice)
	// this should be checked -> case "count": return compare(c.Operator, c.Value, )
	default:
		return false
	}
}

func compare(operator string, left interface{}, right interface{}) bool {
	switch strings.ToUpper(operator) {
	case "=":
		return left == right
	case "!=":
		return left != right
	case ">":
		return left.(float64) > right.(float64)
	case ">=":
		return left.(float64) >= right.(float64)
	case "<":
		return left.(float64) < right.(float64)
	case "<=":
		return left.(float64) <= right.(float64)
	case "IN":
		for _, e := range left.([]string) {
			if e == right.(string) {
				return true
			}
		}
	default:
		return false
	}
	return false
}

// TODO: ADD SEPARATE FUNCTION FOR EACH ACTION TYPE FOR PRODUCT, CATEGORY, AND CART
func ApplyAction(a *parser.Action, store *StoreContext, product *ProductContext) error {
	c := strings.ToLower(a.DiscountType)
	f, err := strconv.ParseFloat(a.Value.(string), 64)
	if err != nil {
		return err
	}
	switch c {
	case "product_percentage":
		product.DiscountedPrice = product.OriginalPrice * (1 - f/100)
		return nil
	case "product_flat_amount":
		product.DiscountedPrice = product.OriginalPrice - f
		return nil
	case "cart_percentage":
		store.CartPriceAfterDiscounts = store.CartPrice * (1 - f/100)
		return nil
	case "cart_flat_amount":
		store.CartPriceAfterDiscounts = store.CartPrice - f
		return nil
	default:
		return fmt.Errorf("unknown discount type %s", c)
	}
}
