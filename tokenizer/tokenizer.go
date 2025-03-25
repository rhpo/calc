package calc

// Token is a token in the input string.

import (
	"fmt"
)

type Token struct {
	// Type is the type of the token.
	Type  TokenType
	Value string
}

// TokenType is the type of a token.
type TokenType int

const (
	NumberToken TokenType = iota
	PlusToken
	MinusToken
	MultiplyToken
	DivideToken
	PowerToken
	EqualToken
	LeftParenthesisToken
	RightParenthesisToken
	Identifier
	Comma
	EOF
)

func (t TokenType) String() string {
	switch t {
	case NumberToken:
		return "Number"
	case PlusToken:
		return "Plus"
	case MinusToken:
		return "Minus"
	case MultiplyToken:
		return "Multiply"
	case DivideToken:
		return "Divide"
	case LeftParenthesisToken:
		return "LeftParenthesis"
	case EqualToken:
		return "EqualToken"
	case PowerToken:
		return "PowerToken"
	case RightParenthesisToken:
		return "RightParenthesis"
	case Identifier:
		return "Identifier"
	case Comma:
		return "Comma"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// Tokenize returns a slice of tokens from the input string.
func Tokenize(input string) []Token {
	var tokens []Token
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case ' ':
			continue

		case ',':
			tokens = append(tokens, Token{Type: Comma, Value: ","})

		case '+':
			tokens = append(tokens, Token{Type: PlusToken, Value: "+"})
		case '-':
			tokens = append(tokens, Token{Type: MinusToken, Value: "-"})
		case '*':

			if i < len(input)-1 && input[i+1] == '*' {
				i++
				tokens = append(tokens, Token{Type: PowerToken, Value: "**"})
				continue
			}

			tokens = append(tokens, Token{Type: MultiplyToken, Value: "*"})
		case '/':
			tokens = append(tokens, Token{Type: DivideToken, Value: "/"})
		case '(':
			tokens = append(tokens, Token{Type: LeftParenthesisToken, Value: "("})
		case ')':
			tokens = append(tokens, Token{Type: RightParenthesisToken, Value: ")"})
		case '\n':
			// do nothing
		default:

			// if it's a number, read the whole number, support floating point numbers, and add it to the tokens
			// if it's an identifier, read the whole identifier, and add it to the tokens

			// if it's a number
			if input[i] >= '0' && input[i] <= '9' {
				start := i
				for i < len(input) && (input[i] >= '0' && input[i] <= '9' || input[i] == '.') {
					i++
				}
				tokens = append(tokens, Token{Type: NumberToken, Value: input[start:i]})
				i--
			} else if (input[i] >= 'a' && input[i] <= 'z') || (input[i] >= 'A' && input[i] <= 'Z') {
				start := i
				for i < len(input) && ((input[i] >= 'a' && input[i] <= 'z') || (input[i] >= 'A' && input[i] <= 'Z')) {
					i++
				}
				tokens = append(tokens, Token{Type: Identifier, Value: input[start:i]})
				i--
			} else {
				panic(fmt.Sprintf("unexpected character: %c", input[i]))
			}
		}
	}

	// tokens = append(tokens, Token{Type: EOF, Value: ""})

	return tokens
}

func PrintTokens(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("%v ", token.Value)
	}
	fmt.Println()
}
