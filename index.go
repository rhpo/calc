package main

import (
	"bufio"
	eval "calc/eval"
	parser "calc/parser"
	tokenizer "calc/tokenizer"

	"fmt"
	"os"
	"strings"
)

const prompt string = ">> "
const message string = "Calc v1.0, type 'exit' to leave."

func main() {

	fmt.Println(message)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(prompt)

		// read input, but avoid using Scanln or Scanf, as it will not read the entire line
		var input string
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input) // Removes trailing newline

		if input == "" {
			continue
		} else if input == "exit" {
			break
		}

		tokens := tokenizer.Tokenize(input)
		// tokenizer.PrintTokens(tokens)

		p := parser.NewParser(tokens)
		program, err := p.Parse()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// parser.PrintProgram(program)
		// fmt.Println(program)

		result, err := eval.Eval(program)
		if err != nil {
			fmt.Println("Evaluation error -", err)
			continue
		}

		fmt.Println(result)
	}
}
