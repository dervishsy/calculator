package parser

import (
	"errors"
	"fmt"
	"strconv"
)

// Token представляет токен в выражении
type Token struct {
	Type  TokenType
	Value string
}

// TokenType представляет тип токена
type TokenType int

const (
	Number TokenType = iota
	Plus
	Minus
	Multiply
	Divide
	LeftParen
	RightParen
	Empty
)

// Node представляет узел в дереве выражения
type Node struct {
	Token  Token
	Left   *Node
	Right  *Node
	Value  float64
	Parsed bool
}

func Parse(expr string) (*Node, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}

	root, remaining, err := parseExpression(tokens, 0)
	if err != nil {
		return nil, err
	}
	if len(remaining) > 0 {
		return nil, errors.New("unexpected token in expression")
	}

	return root, nil
}

func (n *Node) Evaluate() (float64, error) {
	if n.Parsed {
		return n.Value, nil
	}

	left, err := n.Left.Evaluate()
	if err != nil {
		return 0, err
	}

	right, err := n.Right.Evaluate()
	if err != nil {
		return 0, err
	}

	switch n.Token.Type {
	case Plus:
		n.Value = left + right
	case Minus:
		n.Value = left - right
	case Multiply:
		n.Value = left * right
	case Divide:
		if right == 0 {
			return 0, errors.New("division by zero")
		}
		n.Value = left / right
	default:
		return 0, fmt.Errorf("unknown operator: %s", n.Token.Value)
	}

	n.Parsed = true
	return n.Value, nil
}

func parseExpression(tokens []Token, start int) (*Node, []Token, error) {
	left, remaining, err := parseTerm(tokens, start)
	if err != nil {
		return nil, nil, err
	}

	for len(remaining) > 0 && (remaining[0].Type == Plus || remaining[0].Type == Minus) {
		op := remaining[0]
		remaining = remaining[1:]

		right, remaining2, err := parseTerm(remaining, 0)
		if err != nil {
			return nil, nil, err
		}

		left = &Node{op, left, right, 0, false}
		remaining = remaining2
	}

	return left, remaining, nil
}

func tokenize(expr string) ([]Token, error) {
	var tokens []Token
	for len(expr) > 0 {
		token, remainder, err := nextToken(expr)
		if err != nil {
			return nil, err
		}
		if token.Type != Empty {
			tokens = append(tokens, token)
		}
		expr = remainder
	}
	return tokens, nil
}

func nextToken(expr string) (Token, string, error) {
	if len(expr) == 0 {
		return Token{}, "", nil
	}

	switch expr[0] {
	case '+':
		return Token{Plus, "+"}, expr[1:], nil
	case '-':
		return Token{Minus, "-"}, expr[1:], nil
	case '*':
		return Token{Multiply, "*"}, expr[1:], nil
	case '/':
		return Token{Divide, "/"}, expr[1:], nil
	case '(':
		return Token{LeftParen, "("}, expr[1:], nil
	case ')':
		return Token{RightParen, ")"}, expr[1:], nil
	case ' ':
		return Token{Empty, ""}, expr[1:], nil
	case '\t':
		return Token{Empty, ""}, expr[1:], nil
	default:
		if isDigit(expr[0]) {
			num, end := extractNumber(expr)
			return Token{Number, num}, expr[end:], nil
		}
		return Token{}, "", fmt.Errorf("unexpected character: %c", expr[0])
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9' || c == '.'
}

func extractNumber(s string) (string, int) {
	var i int
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	return s[:i], i
}

func parseTerm(tokens []Token, start int) (*Node, []Token, error) {
	left, remaining, err := parseFactor(tokens, start)
	if err != nil {
		return nil, nil, err
	}

	for len(remaining) > 0 && (remaining[0].Type == Multiply || remaining[0].Type == Divide) {
		op := remaining[0]
		remaining = remaining[1:]

		right, remaining2, err := parseFactor(remaining, 0)
		if err != nil {
			return nil, nil, err
		}

		left = &Node{op, left, right, 0, false}
		remaining = remaining2
	}

	return left, remaining, nil
}

func parseFactor(tokens []Token, start int) (*Node, []Token, error) {
	if start >= len(tokens) {
		return nil, nil, errors.New("unexpected end of expression")
	}

	token := tokens[start]
	switch token.Type {
	case Number:
		value, _ := strconv.ParseFloat(token.Value, 64)
		return &Node{token, nil, nil, value, true}, tokens[start+1:], nil
	case Plus:
		right, remaining, err := parseFactor(tokens, start+1)
		if err != nil {
			return nil, nil, err
		}
		return &Node{token, &Node{Token{Number, "0"}, nil, nil, 0, true}, right, 0, false}, remaining, nil
	case Minus:
		right, remaining, err := parseFactor(tokens, start+1)
		if err != nil {
			return nil, nil, err
		}
		negRight := &Node{Token{Minus, "-"}, &Node{Token{Number, "0"}, nil, nil, 0, true}, right, 0, false}
		return negRight, remaining, nil
	case LeftParen:
		expr, remaining, err := parseExpression(tokens[start+1:], 0)
		if err != nil {
			return nil, nil, err
		}
		if len(remaining) == 0 || remaining[0].Type != RightParen {
			return nil, nil, errors.New("missing closing parenthesis")
		}
		return expr, remaining[1:], nil
	default:
		return nil, nil, fmt.Errorf("unexpected token: %s", token.Value)
	}
}
