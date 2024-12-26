package parser

import (
	"fmt"
	"strings"

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

func (p *Parser) parseCondition() (*Condition, error) {
	key := strings.ToUpper(p.currentToken.Value)
	if err := p.consume(lexer.TokenKeyword); err != nil {
		return nil, err
	}

	operator := p.currentToken.Value
	if err := p.consume(lexer.TokenComparison); err != nil {
		return nil, err
	}

	var value interface{}
	switch p.currentToken.Type {
	case lexer.TokenOpenBracket:
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
			if p.currentToken.Type == lexer.TokenLiteralNumber {
				if err := p.consume(lexer.TokenLiteralNumber); err != nil {
					return nil, err
				}
			} else if p.currentToken.Type == lexer.TokenUUID {
				if err := p.consume(lexer.TokenUUID); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Value)
			}
		}
		value = values
		if err := p.consume(lexer.TokenCloseBracket); err != nil {
			return nil, err
		}
	case lexer.TokenLiteralNumber:
		value = p.currentToken.Value
		if err := p.consume(lexer.TokenLiteralNumber); err != nil {
			return nil, err
		}
	case lexer.TokenLiteralString:
		value = p.currentToken.Value
		if err := p.consume(lexer.TokenLiteralString); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Value)
	}

	return &Condition{Key: key, Operator: operator, Value: value}, nil
}

func (p *Parser) parseAction() (*Action, error) {
	if err := p.consume(lexer.TokenAction); err != nil {
		return nil, err
	}

	discountType := p.currentToken.Value
	if err := p.consume(lexer.TokenDiscountType); err != nil {
		return nil, err
	}

	value := p.currentToken.Value
	if err := p.consume(lexer.TokenLiteralNumber); err != nil {
		return nil, err
	}

	return &Action{DiscountType: discountType, Value: value}, nil
}

func (p *Parser) ParseRule() (*Rule, error) {
	condition, err := p.parseCondition()
	if err != nil {
		return nil, err
	}

	logicalCondition := LogicalCondition{
		Left: condition,
	}

	if p.currentToken.Type == lexer.TokenLogical {
		logicalCondition.Operator = p.currentToken.Value
		if err := p.consume(lexer.TokenLogical); err != nil {
			return nil, err
		}
		rightCondition, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		logicalCondition.Right = rightCondition
	}

	action, err := p.parseAction()
	if err != nil {
		return nil, err
	}

	return &Rule{
		Condition: &logicalCondition,
		Action:    action,
	}, nil
}
