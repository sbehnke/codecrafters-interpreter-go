package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	LeftParen uint = iota
	RightParen
	EOF
)

type Token struct {
	TokenType uint
	Token     string
	TokenData string
}

func (t Token) String() string {
	tokenData := t.TokenData
	if len(t.TokenData) == 0 {
		tokenData = "null"
	}
	values := []string{tokenTypeToString(t.TokenType), t.Token, tokenData}
	return strings.Join(values, " ")
}

func tokenTypeToString(tokenType uint) string {
	switch tokenType {
	case LeftParen:
		return "LEFT_PAREN"
	case RightParen:
		return "RIGHT_PAREN"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage
	//
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens := make([]Token, 0)
	if len(fileContents) > 0 {
		for _, c := range fileContents {
			switch c {
			case '(':
				tokens = append(tokens, Token{TokenType: LeftParen, Token: "(", TokenData: ""})
			case ')':
				tokens = append(tokens, Token{TokenType: RightParen, Token: ")", TokenData: ""})
			}
		}
		tokens = append(tokens, Token{TokenType: EOF, Token: "", TokenData: ""})

		for _, t := range tokens {
			fmt.Println(t.String())
		}
	} else {
		fmt.Println("EOF  null") // Placeholder, replace this line when implementing the scanner
	}
}
