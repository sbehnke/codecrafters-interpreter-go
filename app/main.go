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
	Assignment

	Bang
	BangEqual
	Equal
	EqualEqual
	Less
	LessEqual
	Greater
	GreaterEqual

	String

	EOF
)

type Token struct {
	TokenType uint
	Token     string
	TokenData *string
}

type Parser struct {
	Idx0             int
	Idx1             int
	peek             int
	line             int
	HasLexicalErrors bool
	Source           []byte
	Tokens           []Token
}

func NewParser(source []byte) *Parser {
	return &Parser{
		Source:           source,
		Idx0:             0,
		Idx1:             1,
		peek:             0,
		line:             1,
		HasLexicalErrors: false,
		Tokens:           make([]Token, 0),
	}
}

func (p Parser) Chr0() byte {
	if len(p.Source) > p.Idx0 {
		return p.Source[p.Idx0]
	}
	return 0
}

func (p Parser) Chr1() byte {
	if len(p.Source) > p.Idx1 {
		return p.Source[p.Idx1]
	}
	return 0
}

func (p *Parser) Peek() byte {
	p.peek += 1
	if p.peek < len(p.Source) {
		return p.Source[p.peek]
	}

	return 0
}

func (p Parser) Match(c byte) bool {
	return p.Chr1() == c
}

func (p *Parser) Next() byte {
	if p.Idx0 < len(p.Source) {
		p.Idx0 += 1
		p.Idx1 += 1
		p.peek = p.Idx0
	}

	return p.Chr0()
}

func (p Parser) PrintTokens() {
	for _, t := range p.Tokens {
		fmt.Println(t.String())
	}
}

func (p *Parser) Tokenize() {
	p.line = 1
	c := p.Chr0()

	if c != 0 {
		for {
			switch c {
			case '(':
				p.Tokens = append(p.Tokens, Token{TokenType: LeftParen, Token: string(c), TokenData: nil})
			case ')':
				p.Tokens = append(p.Tokens, Token{TokenType: RightParen, Token: string(c), TokenData: nil})
			case '{':
				p.Tokens = append(p.Tokens, Token{TokenType: LeftBrace, Token: string(c), TokenData: nil})
			case '}':
				p.Tokens = append(p.Tokens, Token{TokenType: RightBrace, Token: string(c), TokenData: nil})
			case '*':
				p.Tokens = append(p.Tokens, Token{TokenType: Star, Token: string(c), TokenData: nil})
			case '.':
				p.Tokens = append(p.Tokens, Token{TokenType: Dot, Token: string(c), TokenData: nil})
			case ',':
				p.Tokens = append(p.Tokens, Token{TokenType: Comma, Token: string(c), TokenData: nil})
			case '+':
				p.Tokens = append(p.Tokens, Token{TokenType: Plus, Token: string(c), TokenData: nil})
			case '-':
				p.Tokens = append(p.Tokens, Token{TokenType: Minus, Token: string(c), TokenData: nil})
			case ';':
				p.Tokens = append(p.Tokens, Token{TokenType: Semicolon, Token: string(c), TokenData: nil})

			case '!':
				if p.Match('=') {
					p.Tokens = append(p.Tokens, Token{TokenType: BangEqual, Token: "!=", TokenData: nil})
					p.Next()
				} else {
					p.Tokens = append(p.Tokens, Token{TokenType: Bang, Token: string(c), TokenData: nil})
				}
			case '=':
				if p.Match('=') {
					p.Tokens = append(p.Tokens, Token{TokenType: EqualEqual, Token: "==", TokenData: nil})
					p.Next()
				} else {
					p.Tokens = append(p.Tokens, Token{TokenType: Equal, Token: string(c), TokenData: nil})
				}
			case '>':
				if p.Match('=') {
					p.Tokens = append(p.Tokens, Token{TokenType: GreaterEqual, Token: ">=", TokenData: nil})
					p.Next()
				} else {
					p.Tokens = append(p.Tokens, Token{TokenType: Greater, Token: string(c), TokenData: nil})
				}
			case '<':
				if p.Match('=') {
					p.Tokens = append(p.Tokens, Token{TokenType: LessEqual, Token: "<=", TokenData: nil})
					p.Next()
				} else {
					p.Tokens = append(p.Tokens, Token{TokenType: Less, Token: string(c), TokenData: nil})
				}

			case '/':
				if p.Match('/') {
					for {
						peek := p.Peek()
						if peek == 0 || peek == '\n' {
							break
						} else {
							p.Next()
						}
					}
				} else {
					p.Tokens = append(p.Tokens, Token{TokenType: Slash, Token: string(c), TokenData: nil})
				}

			case '"':
				start := p.Idx0
				for {
					peek := p.Peek()
					if peek == '"' {
						p.Next()
						stop := p.Idx1
						token := string(p.Source[start:stop])
						value := string(p.Source[start+1 : stop-1])
						p.Tokens = append(p.Tokens, Token{TokenType: String, Token: token, TokenData: StrPtr(value)})
						break
					} else if peek == 0 {
						fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", p.line)
						p.HasLexicalErrors = true
						break
					}
					p.Next()
				}

			case ' ':
			case '\t':
			case '\r':
				{
					break
				}

			case '\n':
				p.line += 1

			default:
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", p.line, string(c))
				p.HasLexicalErrors = true
			}

			c = p.Next()
			if c == 0 {
				break
			}
		}
	}
	p.Tokens = append(p.Tokens, Token{TokenType: EOF, Token: "", TokenData: nil})
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
	case Bang:
		return "BANG"
	case BangEqual:
		return "BANG_EQUAL"
	case Equal:
		return "EQUAL"
	case EqualEqual:
		return "EQUAL_EQUAL"
	case Less:
		return "LESS"
	case LessEqual:
		return "LESS_EQUAL"
	case Greater:
		return "GREATER"
	case GreaterEqual:
		return "GREATER_EQUAL"
	case String:
		return "STRING"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if command == "tokenize" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		parser := NewParser(fileContents)
		parser.Tokenize()
		parser.PrintTokens()

		if parser.HasLexicalErrors {
			os.Exit(LexicalError)
		}
	}
}
