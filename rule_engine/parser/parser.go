package parser

import (
	"fmt"

	"github.com/ragnarrlaw/rules/rule_engine/lexer"
)

type Condition struct {
	Key      string
	Operator string
	Value    interface{}
}

type Action struct {
	DiscountType string
	Value        interface{}
}

type LogicalCondition struct {
	Left     *Condition
	Operator string
	Right    *Condition
}

type Rule struct {
	Condition *LogicalCondition
	Action    *Action
}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{lexer: l, currentToken: l.NextToken()}
}

func (p *Parser) consume(expectedType lexer.TokenType) error {
	if p.currentToken.Type != expectedType {
		return fmt.Errorf("unexpected token: %s", p.currentToken.Value)
	}
	p.currentToken = p.lexer.NextToken()
	return nil
}

/*
Parses a given condition such as,
1. product_id in [1,2,3]
2. product_id == 10

Output for the second case:

	Condition{
	  Key: product_id
	  Operator: '=='
	  Value: 10
	}
*/
func (p *Parser) parseCondition() (*Condition, error) {
	key := p.currentToken.Value
	if err := p.consume(lexer.TokenKeyword); err != nil {
		return nil, err
	}

	operator := p.currentToken.Value
	if err := p.consume(lexer.TokenComparison); err != nil {
		return nil, err
	}

	var value interface{}
	if p.currentToken.Type == lexer.TokenOpenBracket { // for [1,2,3]
		if err := p.consume(lexer.TokenOpenBracket); err != nil {
			return nil, err
		}
		var values []string
		for p.currentToken.Type != lexer.TokenCloseBracket {
			if p.currentToken.Type == lexer.TokenComma {
				if err := p.consume(lexer.TokenComma); err != nil {
					return nil, err
				}
			}
			values = append(values, p.currentToken.Value)
			if err := p.consume(lexer.TokenLiteralNumber); err != nil {
				return nil, err
			}
		}
		value = values
		if err := p.consume(lexer.TokenCloseBracket); err != nil {
			return nil, err
		}
	} else {
		value = p.currentToken.Value
		if err := p.consume(lexer.TokenLiteralNumber); err != nil {
			return nil, err
		}
	}
	return &Condition{Key: key, Operator: operator, Value: value}, nil
}

/*
Parses a given action such as,
1.THEN percentage = 10
2.THEN flat = 5
3.THEN bogo = 1

Output for the second case:

	Action{
	  DiscountType: percentage
	  Value: 10
	}
*/
func (p *Parser) parseAction() (*Action, error) {
	if err := p.consume(lexer.TokenAction); err != nil {
		return nil, err
	}

	discountType := p.currentToken.Value
	if err := p.consume(lexer.TokenDiscountType); err != nil {
		return nil, err
	}

	if err := p.consume(lexer.TokenAssignmentOperator); err != nil {
		return nil, err
	}

	value := p.currentToken.Value
	if err := p.consume(lexer.TokenLiteralNumber); err != nil {
		return nil, err
	}

	return &Action{DiscountType: discountType, Value: value}, nil
}

/*
Parses a single rule with two conditions connected by a single operator
*/
func (p *Parser) ParseRule() (*Rule, error) {
	condition, err := p.parseCondition()
	if err != nil {
		return nil, err
	}

	LogicalCondition := LogicalCondition{
		Left: condition,
	}

	if p.currentToken.Type == lexer.TokenLogical {
		LogicalCondition.Operator = p.currentToken.Value
		if err := p.consume(lexer.TokenLogical); err != nil {
			return nil, err
		}
		rightCondition, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		LogicalCondition.Right = rightCondition
	}

	action, err := p.parseAction()
	if err != nil {
		return nil, err
	}

	return &Rule{
		Condition: &LogicalCondition,
		Action:    action,
	}, nil
}
