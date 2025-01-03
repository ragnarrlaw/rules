package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []*Token
	}{
		{
			`product_id IN [1, 2, 3]`,
			[]*Token{
				{Type: TokenKeyword, Value: "PRODUCT_ID"},
				{Type: TokenComparison, Value: "IN"},
				{Type: TokenOpenBracket, Value: "["},
				{Type: TokenLiteralNumber, Value: "1"},
				{Type: TokenComma, Value: ","},
				{Type: TokenLiteralNumber, Value: "2"},
				{Type: TokenComma, Value: ","},
				{Type: TokenLiteralNumber, Value: "3"},
				{Type: TokenCloseBracket, Value: "]"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			`total_price > 100 AND total_category_price < 500`,
			[]*Token{
				{Type: TokenKeyword, Value: "TOTAL_PRICE"},
				{Type: TokenComparison, Value: ">"},
				{Type: TokenLiteralNumber, Value: "100"},
				{Type: TokenLogical, Value: "AND"},
				{Type: TokenKeyword, Value: "TOTAL_CATEGORY_PRICE"},
				{Type: TokenComparison, Value: "<"},
				{Type: TokenLiteralNumber, Value: "500"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			`THEN Percentage 10`,
			[]*Token{
				{Type: TokenAction, Value: "THEN"},
				{Type: TokenDiscountType, Value: "PERCENTAGE"},
				{Type: TokenLiteralNumber, Value: "10"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			`product_id IN ["9f9285c6-a4d3-407e-9bd6-92ed094d0b02"]`,
			[]*Token{
				{Type: TokenKeyword, Value: "PRODUCT_ID"},
				{Type: TokenComparison, Value: "IN"},
				{Type: TokenOpenBracket, Value: "["},
				{Type: TokenUUID, Value: "9f9285c6-a4d3-407e-9bd6-92ed094d0b02"},
				{Type: TokenCloseBracket, Value: "]"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			`product_id = "9f9285c6-a4d3-407e-9bd6-92ed094d0b02"`,
			[]*Token{
				{Type: TokenKeyword, Value: "PRODUCT_ID"},
				{Type: TokenComparison, Value: "="},
				{Type: TokenUUID, Value: "9f9285c6-a4d3-407e-9bd6-92ed094d0b02"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, test := range tests {
		l := NewLexer(test.input)
		for i, expected := range test.expected {
			token := l.NextToken()
			if token.Type != expected.Type || token.Value != expected.Value {
				t.Fatalf(
					"test failed for input '%s'. token %d: expected %v, got %v",
					test.input,
					i,
					expected,
					token,
				)
			}
		}
	}
}
