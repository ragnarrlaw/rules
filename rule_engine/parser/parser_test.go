package parser

import (
	"testing"

	"github.com/ragnarrlaw/rules/rule_engine/lexer"
)

func TestParseRule(t *testing.T) {
	tests := []struct {
		input    string
		expected Rule
	}{
		{
			`product_id in [1, 2, 3] AND product_id in [1, 2, 3] then percentage 10`,
			Rule{
				Condition: &LogicalCondition{
					Left: &Condition{
						Key:      "PRODUCT_ID",
						Operator: "IN",
						Value:    []string{"1", "2", "3"},
					},
					Operator: "AND",
					Right: &Condition{
						Key:      "PRODUCT_ID",
						Operator: "IN",
						Value:    []string{"1", "2", "3"},
					},
				}, Action: &Action{
					DiscountType: "PERCENTAGE",
					Value:        "10",
				},
			},
		},
		{
			`product_id in ["9f9285c6-a4d3-407e-9bd6-92ed094d0b02"] then percentage 10`,
			Rule{
				Condition: &LogicalCondition{
					Left: &Condition{
						Key:      "PRODUCT_ID",
						Operator: "IN",
						Value:    []string{`"9f9285c6-a4d3-407e-9bd6-92ed094d0b02"`},
					},
				}, Action: &Action{
					DiscountType: "PERCENTAGE",
					Value:        "10",
				},
			},
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		r, err := p.ParseRule()
		if err != nil {
			t.Fatalf(
				"parsing rule %s with the error %s",
				test.input,
				err.Error(),
			)
		}
		if r.Condition.Operator != test.expected.Condition.Operator {
			t.Fatalf(
				"operator %s interpreted as %s",
				test.expected.Condition.Operator,
				r.Condition.Operator,
			)
		}
		if r.Condition.Left.Operator != test.expected.Condition.Left.Operator {
			t.Fatalf(
				"left condition operator %s interpreted as %s",
				test.expected.Condition.Left.Operator,
				r.Condition.Left.Operator,
			)
		}
		if r.Condition.Left.Key != test.expected.Condition.Left.Key {
			t.Fatalf(
				"left condition key %s interpreted as %s",
				test.expected.Condition.Left.Key,
				r.Condition.Left.Key,
			)
		}
		if r.Condition.Left.Value != nil {
			for i, v := range r.Condition.Left.Value.([]string) {
				if v != test.expected.Condition.Left.Value.([]string)[i] {
					t.Fatalf(
						"left condition value at %d index %s, interpreted as %s",
						i,
						test.expected.Condition.Left.Value.([]string)[i],
						v,
					)
				}
			}
		}
		if test.expected.Condition.Operator != "" {
			if r.Condition.Right.Operator != test.expected.Condition.Right.Operator {
				t.Fatalf(
					"right condition operator %s interpreted as %s",
					test.expected.Condition.Right.Operator,
					r.Condition.Right.Operator,
				)
			}
			if r.Condition.Right.Key != test.expected.Condition.Right.Key {
				t.Fatalf(
					"right condition key %s interpreted as %s",
					test.expected.Condition.Right.Key,
					r.Condition.Right.Key,
				)
			}
		}
	}
}
