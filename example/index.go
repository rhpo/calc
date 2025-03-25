package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TokenType int

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	LPAREN
	RPAREN
	FUNCTION
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	input string
	pos   int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: strings.ReplaceAll(input, " ", ""), pos: 0}
}

func (t *Tokenizer) getNextToken() Token {
	if t.pos >= len(t.input) {
		return Token{Type: EOF, Value: ""}
	}

	ch := t.input[t.pos]

	if ch >= '0' && ch <= '9' || ch == '.' {
		start := t.pos
		for t.pos < len(t.input) && (t.input[t.pos] >= '0' && t.input[t.pos] <= '9' || t.input[t.pos] == '.') {
			t.pos++
		}
		return Token{Type: NUMBER, Value: t.input[start:t.pos]}
	}

	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
		start := t.pos
		for t.pos < len(t.input) && ((t.input[t.pos] >= 'a' && t.input[t.pos] <= 'z') || (t.input[t.pos] >= 'A' && t.input[t.pos] <= 'Z')) {
			t.pos++
		}
		return Token{Type: FUNCTION, Value: t.input[start:t.pos]}
	}

	t.pos++
	switch ch {
	case '+':
		return Token{Type: PLUS, Value: string(ch)}
	case '-':
		return Token{Type: MINUS, Value: string(ch)}
	case '*':
		return Token{Type: MULTIPLY, Value: string(ch)}
	case '/':
		return Token{Type: DIVIDE, Value: string(ch)}
	case '(':
		return Token{Type: LPAREN, Value: string(ch)}
	case ')':
		return Token{Type: RPAREN, Value: string(ch)}
	default:
		panic("Unknown character: " + string(ch))
	}
}

type Node interface {
	Evaluate() float64
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) Evaluate() float64 {
	return n.Value
}

type BinaryOpNode struct {
	Left     Node
	Operator TokenType
	Right    Node
}

func (b *BinaryOpNode) Evaluate() float64 {
	switch b.Operator {
	case PLUS:
		return b.Left.Evaluate() + b.Right.Evaluate()
	case MINUS:
		return b.Left.Evaluate() - b.Right.Evaluate()
	case MULTIPLY:
		return b.Left.Evaluate() * b.Right.Evaluate()
	case DIVIDE:
		return b.Left.Evaluate() / b.Right.Evaluate()
	default:
		panic("Unknown operator")
	}
}

type FunctionNode struct {
	Function string
	Arg      Node
}

func (f *FunctionNode) Evaluate() float64 {
	switch f.Function {
	case "sin":
		return math.Sin(f.Arg.Evaluate())
	case "cos":
		return math.Cos(f.Arg.Evaluate())
	default:
		panic("Unknown function: " + f.Function)
	}
}

type Parser struct {
	tokenizer *Tokenizer
	current   Token
}

func NewParser(tokenizer *Tokenizer) *Parser {
	parser := &Parser{tokenizer: tokenizer}
	parser.current = tokenizer.getNextToken()
	return parser
}

func (p *Parser) eat(tokenType TokenType) {
	if p.current.Type == tokenType {
		p.current = p.tokenizer.getNextToken()
	} else {
		panic("Unexpected token")
	}
}

func (p *Parser) parseFactor() Node {
	if p.current.Type == NUMBER {
		val, _ := strconv.ParseFloat(p.current.Value, 64)
		node := &NumberNode{Value: val}
		p.eat(NUMBER)
		return node
	}
	if p.current.Type == FUNCTION {
		funcName := p.current.Value
		p.eat(FUNCTION)
		p.eat(LPAREN)
		arg := p.parseExpression()
		p.eat(RPAREN)
		return &FunctionNode{Function: funcName, Arg: arg}
	}
	if p.current.Type == LPAREN {
		p.eat(LPAREN)
		node := p.parseExpression()
		p.eat(RPAREN)
		return node
	}
	panic("Unexpected token")
}

func (p *Parser) parseTerm() Node {
	node := p.parseFactor()
	for p.current.Type == MULTIPLY || p.current.Type == DIVIDE {
		op := p.current.Type
		p.eat(op)
		node = &BinaryOpNode{Left: node, Operator: op, Right: p.parseFactor()}
	}
	return node
}

func (p *Parser) parseExpression() Node {
	node := p.parseTerm()
	for p.current.Type == PLUS || p.current.Type == MINUS {
		op := p.current.Type
		p.eat(op)
		node = &BinaryOpNode{Left: node, Operator: op, Right: p.parseTerm()}
	}
	return node
}

func main() {
	fmt.Print("Enter expression: ")
	var input string
	fmt.Scanln(&input)
	tokenizer := NewTokenizer(input)
	parser := NewParser(tokenizer)
	tree := parser.parseExpression()

	fmt.Println("Result:", tree.Evaluate())
}
