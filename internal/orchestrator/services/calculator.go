package services

import (
	"errors"
	"strconv"
	"strings"
)

type parser struct {
	expression string
	pos        int
	current    rune
}

func (p *parser) next() {
	p.pos++
	if p.pos < len(p.expression) {
		p.current = rune(p.expression[p.pos])
	} else {
		p.current = 0
	}
}

func (p *parser) parse() (float64, error) {
	p.next() 
	result, err := p.parseExpression()
	if err != nil {
		return 0, err
	}
	if p.current != 0 {
		return 0, errors.New("unexpected character")
	}
	return result, nil
}

func (p *parser) parseExpression() (float64, error) {
	result, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for p.current == '+' || p.current == '-' {
		op := p.current
		p.next()
		nextTerm, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if op == '+' {
			result += nextTerm
		} else {
			result -= nextTerm
		}
	}
	return result, nil
}

func (p *parser) parseTerm() (float64, error) {
	result, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for p.current == '*' || p.current == '/' {
		op := p.current
		p.next()
		nextFactor, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		if op == '*' {
			result *= nextFactor
		} else {
			if nextFactor == 0 {
				return 0, errors.New("division by zero")
			}
			result /= nextFactor
		}
	}
	return result, nil
}

func (p *parser) parseFactor() (float64, error) {
	if p.current == '(' {
		p.next()
		result, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if p.current != ')' {
			return 0, errors.New("mismatched parentheses")
		}
		p.next()
		return result, nil
	}

	startPos := p.pos
	for (p.current >= '0' && p.current <= '9') || p.current == '.' {
		p.next()
	}
	if startPos == p.pos {
		return 0, errors.New("expected number")
	}
	value, err := strconv.ParseFloat(p.expression[startPos:p.pos], 64)
	if err != nil {
		return 0, errors.New("invalid number")
	}
	return value, nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	p := parser{expression: expression, pos: -1}
	return p.parse()
}

func EvaluateExpression(expr string) (float64, error) {
	return Calc(expr)
}
