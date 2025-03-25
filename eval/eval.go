package calc

import (
	parser "calc/parser"
	tokenizer "calc/tokenizer"
	"fmt"
	. "math"
	"strconv"
)

func evalBinaryExpression(e *parser.BinaryExpression) (float64, error) {
	left, err := evalExpression(e.Left)
	right, err2 := evalExpression(e.Right)

	if err != nil {
		return 0, err
	} else if err2 != nil {
		return 0, err2
	}

	switch e.Op.Type {
	case tokenizer.PlusToken:
		return left + right, nil
	case tokenizer.MinusToken:
		return left - right, nil
	case tokenizer.MultiplyToken:
		return left * right, nil
	case tokenizer.DivideToken:
		return left / right, nil
	case tokenizer.PowerToken:
		return Pow(left, right), nil
	default:
		return 0, fmt.Errorf("Unknown binary operator: %s", e.Op.Type)
	}
}

// math functions map
var mathFuncs = map[string]func(float64) float64{
	"sin":   Sin,
	"cos":   Cos,
	"tan":   Tan,
	"asin":  Asin,
	"acos":  Acos,
	"atan":  Atan,
	"sqrt":  Sqrt,
	"log":   Log,
	"exp":   Exp,
	"abs":   Abs,
	"ceil":  Ceil,
	"floor": Floor,
	"round": Round,
	"trunc": Trunc,
}

var MathConstants = map[string]float64{
	"pi":      Pi,
	"e":       E,
	"testvar": 1,
}

func evalIdentifierExpression(e *parser.IdentifierExpression) (float64, error) {

	if C, ok := MathConstants[e.Value.Value]; ok {
		return C, nil
	}

	return 0, fmt.Errorf("Undefined variable %s", e.Value.Value)

}

func evalFunctionExpression(e *parser.FuncExpression) (float64, error) {
	args := make([]float64, len(e.Args))
	for i, arg := range e.Args {
		args[i], _ = evalExpression(arg)
	}

	if f, ok := mathFuncs[e.Name]; ok {
		return f(args[0]), nil
	}

	return 0, fmt.Errorf("Unknown function: %s", e.Name)

}

func evalExpression(e parser.Expression) (float64, error) {
	switch e := e.(type) {
	case parser.NumberExpression:
		value, err := strconv.ParseFloat(e.Value.Value, 64)
		if err != nil {
			return 0, err
		}
		return value, nil

	case parser.FuncExpression:
		res, err := evalFunctionExpression(&e)
		return res, err

	case parser.BinaryExpression:
		res, err := evalBinaryExpression(&e)
		return res, err

	case parser.IdentifierExpression:
		res, err := evalIdentifierExpression(&e)

		return res, err

	case parser.Program:
		res, err := evalExpression(e.Statements[0])
		return res, err

	default:
		return 0, fmt.Errorf("Unknown expression type: %T", e)
	}
}

func Eval(program parser.Program) (float64, error) {
	return evalExpression(program)
}

// Eval evaluates the input string and returns the result.
