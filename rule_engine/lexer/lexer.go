package lexer

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenEOF                TokenType = "eof"
	TokenKeyword            TokenType = "keyword"
	TokenComparison         TokenType = "comparison_operator"
	TokenLogical            TokenType = "logical_operator"
	TokenLiteralString      TokenType = "literal_string"
	TokenLiteralNumber      TokenType = "literal_number"
	TokenAssignmentOperator TokenType = "assignment_operator"
	TokenOpenBracket        TokenType = "open_bracket"
	TokenCloseBracket       TokenType = "close_bracket"
	TokenOpenParenthesis    TokenType = "open_parenthesis"
	TokenCloseParenthesis   TokenType = "close_parenthesis"
	TokenComma              TokenType = "comma"
	TokenAction             TokenType = "action"
	TokenDiscountType       TokenType = "discount_types"
	TokenUUID               TokenType = "uuid"
)

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
	case '=':
		return &Token{Type: TokenComparison, Value: "="}
	case '>':
		return &Token{Type: TokenComparison, Value: ">"}
	case '<':
		return &Token{Type: TokenComparison, Value: "<"}
	default:
		t := l.scanner.TokenText()

		// Check if the token is a UUID
		if _, err := uuid.Parse(t); err == nil {
			return &Token{Type: TokenUUID, Value: t}
		}

		t = strings.ToUpper(t)

		switch t {
		case "==", "!=", ">=", "<=", "IN":
			return &Token{Type: TokenComparison, Value: t}
		}

		if unicode.IsDigit([]rune(t)[0]) || t[0] == '-' || strings.Contains(t, ".") {
			return &Token{Type: TokenLiteralNumber, Value: t}
		}

		switch t {
		case "AND", "OR":
			return &Token{Type: TokenLogical, Value: t}
		case "THEN":
			return &Token{Type: TokenAction, Value: t}
		case "PERCENTAGE", "FLAT_AMOUNT", "BOGO", "PRODUCT_PERCENTAGE", "CART_PERCENTAGE", "PRODUCT_FLAT_AMOUNT", "CART_FLAT_AMOUNT":
			return &Token{Type: TokenDiscountType, Value: t}
		case "MIN_CART_PRICE", "TOTAL_PRICE", "PRODUCT_ID", "TOTAL_CATEGORY_PRICE", "CATEGORY_ID", "PURCHASE_QUANTITY":
			return &Token{Type: TokenKeyword, Value: t}
		default:
			if unicode.IsLetter([]rune(t)[0]) {
				return &Token{Type: TokenLiteralString, Value: t}
			}
			return &Token{Type: TokenEOF, Value: fmt.Sprintf("unexpected character: %s", t)}
		}
	}
}
