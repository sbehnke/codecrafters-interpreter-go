package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	LexicalError = 65
)

const (
	LeftParen uint = iota
	RightParen
	LeftBrace
	RightBrace
	Star
	Dot
	Comma
	Plus
	Minus
	Slash
	Semicolon
	EOF
)

type Token struct {
	TokenType uint
	Token     string
	TokenData *string
}

func StrPtr(s string) *string {
	return &s
}

func (t Token) String() string {
	var tokenData string
	if t.TokenData == nil {
		tokenData = "null"
	} else {
		tokenData = *t.TokenData
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
	case LeftBrace:
		return "LEFT_BRACE"
	case RightBrace:
		return "RIGHT_BRACE"
	case Star:
		return "STAR"
	case Dot:
		return "DOT"
	case Comma:
		return "COMMA"
	case Plus:
		return "PLUS"
	case Minus:
		return "MINUS"
	case Slash:
		return "SLASH"
	case Semicolon:
		return "SEMICOLON"
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

	hasLexicalErrors := false
	tokens := make([]Token, 0)
	if len(fileContents) > 0 {
		line := 1

		for _, c := range fileContents {
			switch c {
			case '(':
				tokens = append(tokens, Token{TokenType: LeftParen, Token: string(c), TokenData: nil})
			case ')':
				tokens = append(tokens, Token{TokenType: RightParen, Token: string(c), TokenData: nil})
			case '{':
				tokens = append(tokens, Token{TokenType: LeftBrace, Token: string(c), TokenData: nil})
			case '}':
				tokens = append(tokens, Token{TokenType: RightBrace, Token: string(c), TokenData: nil})
			case '*':
				tokens = append(tokens, Token{TokenType: Star, Token: string(c), TokenData: nil})
			case '.':
				tokens = append(tokens, Token{TokenType: Dot, Token: string(c), TokenData: nil})
			case ',':
				tokens = append(tokens, Token{TokenType: Comma, Token: string(c), TokenData: nil})
			case '+':
				tokens = append(tokens, Token{TokenType: Plus, Token: string(c), TokenData: nil})
			case '-':
				tokens = append(tokens, Token{TokenType: Minus, Token: string(c), TokenData: nil})
			case '/':
				tokens = append(tokens, Token{TokenType: Slash, Token: string(c), TokenData: nil})
			case ';':
				tokens = append(tokens, Token{TokenType: Semicolon, Token: string(c), TokenData: nil})
			case '\n':
				line += 1
			default:
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, string(c))
				hasLexicalErrors = true
			}
		}
		tokens = append(tokens, Token{TokenType: EOF, Token: "", TokenData: nil})

		for _, t := range tokens {
			fmt.Println(t.String())
		}
	} else {
		fmt.Println("EOF  null") // Placeholder, replace this line when implementing the scanner
	}

	if hasLexicalErrors {
		os.Exit(LexicalError)
	}
}
