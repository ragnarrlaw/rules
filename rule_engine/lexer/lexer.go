package lexer

import (
	"strings"
	"text/scanner"
	"unicode"
)

type TokenType string

const (
	// end-of-file marker
	TokenEOF TokenType = "eof"
	// any permitted keywords min_cart_value, product_category, etc.
	TokenKeyword TokenType = "keyword"
	// comparison operators ==, !=, >, etc.
	TokenComparison TokenType = "comparison_operator"
	// logical operators AND, and OR
	TokenLogical TokenType = "logical_operator"
	// strings or character sequences
	TokenLiteralString TokenType = "literal_string"
	// numbers 1,2,3,...
	TokenLiteralNumber TokenType = "literal_number"
	// '='
	TokenAssignmentOperator TokenType = "assignment_operator"
	// '['
	TokenOpenBracket TokenType = "open_bracket"
	// ']'
	TokenCloseBracket TokenType = "close_bracket"
	// '('
	TokenOpenParenthesis TokenType = "open_parenthesis"
	// ')'
	TokenCloseParenthesis TokenType = "close_parenthesis"
	// ','
	TokenComma TokenType = "comma"
	// marks the start of the action -> 'THEN' (or 'then')
	TokenAction TokenType = "action"
	// keywords that define discount types -> precentage, flat_amount, and bogo (buy one get one free)
	TokenDiscountType TokenType = "discount_types"
)

// represents a single token
type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	scanner scanner.Scanner
}

func NewLexer(input string) *Lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Whitespace ^= 1<<'\n' | 1<<'\t'
	return &Lexer{s}
}

func (l *Lexer) NextToken() *Token {
	r := l.scanner.Scan()
	switch r {
	case scanner.EOF:
		return &Token{Type: TokenEOF}
	case '(':
		return &Token{Type: TokenOpenParenthesis, Value: "("}
	case ')':
		return &Token{Type: TokenCloseParenthesis, Value: ")"}
	case '[':
		return &Token{Type: TokenOpenBracket, Value: "["}
	case ']':
		return &Token{Type: TokenCloseBracket, Value: "]"}
	case ',':
		return &Token{Type: TokenComma, Value: ","}
	default:
		t := l.scanner.TokenText()
		if unicode.IsDigit([]rune(t)[0]) || t[0] == '-' {
			return &Token{Type: TokenLiteralNumber, Value: t}
		}
		switch t {
		case "AND", "OR", "and", "or":
			return &Token{Type: TokenLogical, Value: t}
		case "THEN", "then":
			return &Token{Type: TokenAction, Value: t}
		case "percentage", "Percentage", "flat_amount", "Flat_amount", "bogo", "BOGO":
			return &Token{Type: TokenDiscountType, Value: t}
		case "min_cart_price", "total_price", "Total_price", "product_id", "Product_id", "total_category_price", "category_id":
			return &Token{Type: TokenKeyword, Value: t}
		case "==", "!=", ">", ">=", "<", "<=", "in", "IN":
			return &Token{Type: TokenComparison, Value: t}
		case "=":
			return &Token{Type: TokenAssignmentOperator, Value: t}
		default:
			return &Token{Type: TokenLiteralString, Value: t}
		}
	}
}
