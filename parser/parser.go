package calc

import (
	tokenizer "calc/tokenizer"
	"fmt"
)

type Statement interface{}
type Expression interface{}

type BinaryExpression struct {
	Left  Expression
	Right Expression
	Op    tokenizer.Token
	Kind  string
}

type UnaryExpression struct {
	Right Expression
	Op    tokenizer.Token
	Kind  string
}

type NumberExpression struct {
	Value tokenizer.Token
	Kind  string
}

type IdentifierExpression struct {
	Value tokenizer.Token
	Kind  string
}

type FuncExpression struct {
	Name string
	Args []Expression
	Kind string
}

type Program struct {
	Statements []Statement
	Kind       string
}

type Parser struct {
	tokens  []tokenizer.Token
	current int
}

func NewParser(tokens []tokenizer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Get(pos ...int) *tokenizer.Token {
	position := p.current
	if len(pos) > 0 {
		position = pos[0]
	}
	if position < 0 || position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[position]
}

func (p *Parser) Expect(T tokenizer.TokenType) (tokenizer.Token, error) {
	next := p.Get(p.current + 1)
	if next == nil || next.Type != T {
		return tokenizer.Token{}, fmt.Errorf("expected %s but got %s", T.String(), p.Get().Type.String())
	}
	p.current++
	return *next, nil
}

func (p *Parser) Parse() (Program, error) {
	var statements []Statement
	for p.current < len(p.tokens) {
		stmt, err := p.parseStatement()
		if err != nil {
			return Program{}, err
		}
		statements = append(statements, stmt)
	}
	return Program{Statements: statements, Kind: "Program"}, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAddition()
}

func (p *Parser) parseAddition() (Expression, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}

	for {
		token := p.Get()
		if token == nil || (token.Type != tokenizer.PlusToken && token.Type != tokenizer.MinusToken) {
			break
		}
		op := *token
		p.current++
		right, err := p.parseMultiplication()
		if err != nil {
			return nil, err
		}
		left = BinaryExpression{Left: left, Right: right, Op: op, Kind: "BinaryExpression"}
	}
	return left, nil
}

func (p *Parser) parseMultiplication() (Expression, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}

	for {
		token := p.Get()
		if token == nil || (token.Type != tokenizer.MultiplyToken && token.Type != tokenizer.DivideToken) {
			break
		}
		op := *token
		p.current++
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = BinaryExpression{Left: left, Right: right, Op: op, Kind: "BinaryExpression"}
	}
	return left, nil
}

func (p *Parser) parsePower() (Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		token := p.Get()
		if token == nil || (token.Type != tokenizer.PowerToken && token.Type != tokenizer.DivideToken) {
			break
		}
		op := *token
		p.current++
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = BinaryExpression{Left: left, Right: right, Op: op, Kind: "BinaryExpression"}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	token := p.Get()
	if token == nil {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch token.Type {
	case tokenizer.NumberToken:
		return p.parseNumber()
	case tokenizer.Identifier:

		if p.Get(p.current+1) == nil ||
			(p.Get(p.current+1) != nil &&
				p.Get(p.current+1).Type != tokenizer.LeftParenthesisToken) {
			return p.parseIdentifier()
		}

		if p.Get(p.current+1) != nil && p.Get(p.current+1).Type == tokenizer.LeftParenthesisToken {
			return p.parseFunction()
		}
		return nil, fmt.Errorf("unexpected token: %s", token.Type.String())
	case tokenizer.LeftParenthesisToken:
		return p.parseGroup()
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Type.String())
	}
}

func (p *Parser) parseNumber() (Expression, error) {
	tok := *p.Get()
	p.current++
	return NumberExpression{Value: tok, Kind: "NumberExpression"}, nil
}

func (p *Parser) parseIdentifier() (Expression, error) {
	tok := *p.Get()
	p.current++
	return IdentifierExpression{Kind: "IdentifierExpression", Value: tok}, nil
}

func (p *Parser) parseGroup() (Expression, error) {
	p.current++
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if token := p.Get(); token == nil || token.Type != tokenizer.RightParenthesisToken {
		return nil, fmt.Errorf("expected RightParenthesis but got %s", token.Type.String())
	}
	p.current++
	return expr, nil
}

func (p *Parser) parseFunction() (Expression, error) {
	name := p.Get().Value
	p.current++
	if token := p.Get(); token == nil || token.Type != tokenizer.LeftParenthesisToken {
		return nil, fmt.Errorf("expected LeftParenthesis but got %s", token.Type.String())
	}
	p.current++

	if p.Get().Type == tokenizer.RightParenthesisToken {
		p.current++
		return FuncExpression{Name: name, Args: []Expression{}, Kind: "FuncExpression"}, nil
	}

	var args []Expression
	for {
		if token := p.Get(); token == nil || token.Type == tokenizer.RightParenthesisToken {
			break
		}
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if token := p.Get(); token != nil && token.Type == tokenizer.Comma {
			p.current++
		}
	}
	
	if token := p.Get(); token == nil || token.Type != tokenizer.RightParenthesisToken {
		return nil, fmt.Errorf("expected RightParenthesis but got %s", token.Type.String())
	}
	p.current++
	return FuncExpression{Name: name, Args: args, Kind: "FuncExpression"}, nil
}

func PrintProgram(program Program) {
	for _, stmt := range program.Statements {
		PrintStatement(stmt)
	}
}

func PrintStatement(stmt Statement) {
	switch stmt := stmt.(type) {
	case Expression:
		PrintExpression(stmt)
	}
}

func PrintExpression(expr Expression) {
	switch expr := expr.(type) {
	case BinaryExpression:
		fmt.Print("(")
		PrintExpression(expr.Left)
		PrintToken(expr.Op)
		PrintExpression(expr.Right)
		fmt.Print(")")
	case UnaryExpression:
		PrintToken(expr.Op)
		PrintExpression(expr.Right)
	case NumberExpression:
		PrintToken(expr.Value)
	case FuncExpression:
		fmt.Printf("%s(", expr.Name)
		for i, arg := range expr.Args {
			PrintExpression(arg)
			if i < len(expr.Args)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print(")")
	default:
		panic("unknown expression type")
	}
}

func PrintToken(tok tokenizer.Token) {
	fmt.Print(tok.Value)
}
