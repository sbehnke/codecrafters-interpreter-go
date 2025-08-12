package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	LexicalError   = 65
	SupportUnicode = true
)

type TokenType uint

const (
	// Single-character tokens
	LeftParen TokenType = iota
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

	// One or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Less
	LessEqual
	Greater
	GreaterEqual

	// Literals
	String
	Number
	Identifier

	// Keywords
	And
	Class
	Else
	False
	For
	Fun
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	// Special
	EOF
)

type Token struct {
	TokenType TokenType
	Token     string
	TokenData *string
	Line      int
}

type Lexer struct {
	idx0             int
	idx1             int
	line             int
	column           int
	HasLexicalErrors bool
	Source           []rune
	Tokens           []Token
	keywords         map[string]TokenType
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		Source:           []rune(source),
		idx0:             0,
		idx1:             1,
		line:             1,
		HasLexicalErrors: false,
		Tokens:           make([]Token, 0),
		keywords: map[string]TokenType{
			"and":    And,
			"class":  Class,
			"else":   Else,
			"false":  False,
			"for":    For,
			"fun":    Fun,
			"if":     If,
			"nil":    Nil,
			"or":     Or,
			"print":  Print,
			"return": Return,
			"super":  Super,
			"this":   This,
			"true":   True,
			"var":    Var,
			"while":  While,
		},
	}
}

func (p Lexer) chr0() rune {
	if len(p.Source) > p.idx0 {
		return p.Source[p.idx0]
	}
	return 0
}

func (p Lexer) chr1() rune {
	if len(p.Source) > p.idx1 {
		return p.Source[p.idx1]
	}
	return 0
}

func (p *Lexer) peek() rune {
	if p.idx1 < len(p.Source) {
		return p.Source[p.idx1]
	}
	return 0
}

func (p *Lexer) peekN(n int) rune {
	idx := p.idx0 + n
	if idx < len(p.Source) {
		return p.Source[idx]
	}
	return 0
}

func (p Lexer) match(c rune) bool {
	return p.chr1() == c
}

func (p *Lexer) next() rune {
	if p.idx0 < len(p.Source) {
		if p.chr0() == '\n' {
			p.line++
			p.column = 0
		} else {
			p.column++
		}
		p.idx0++
		p.idx1++
	}
	return p.chr0()
}

func (p Lexer) PrintTokens() {
	for _, t := range p.Tokens {
		fmt.Println(t.String())
	}
}

func (p *Lexer) lexString() (Token, error) {
	start := p.idx0
	for {
		peek := p.peek()
		switch peek {
		case '"':
			p.next()
			stop := p.idx1
			token := string(p.Source[start:stop])
			value := string(p.Source[start+1 : stop-1])
			return Token{TokenType: String, Token: token, TokenData: StrPtr(value), Line: p.line}, nil
		case 0:
			p.HasLexicalErrors = true
			return Token{Line: p.line}, fmt.Errorf("[line %d] Error: Unterminated string", p.line)
		}
		p.next()
	}
}

func (p *Lexer) lexNumber() (Token, error) {
	start := p.idx0
	for {
		peek := p.peek()
		if unicode.IsNumber(peek) || peek == '.' {
			p.next()
		} else {
			stop := p.idx1
			token := string(p.Source[start:stop])
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				p.HasLexicalErrors = true
				return Token{Line: p.line}, fmt.Errorf("[line %d] Error: Invalid number", p.line)
			}
			var f string
			if value == math.Trunc(value) {
				f = strconv.FormatFloat(value, 'f', 1, 64)
			} else {
				f = strconv.FormatFloat(value, 'g', -1, 64)
			}
			return Token{TokenType: Number, Token: token, TokenData: StrPtr(f), Line: p.line}, nil
		}
	}
}

func (p *Lexer) lexIdentifer() Token {
	start := p.idx0
	for {
		peek := p.peek()
		// Check for Zero-Width Joiner (U+200D) and other joining characters
		if unicode.IsLetter(peek) ||
			unicode.IsNumber(peek) ||
			peek == '_' {
			p.next()
		} else if SupportUnicode && !unicode.IsSpace(peek) && !unicode.IsPunct(peek) && (unicode.IsLetter(peek) ||
			unicode.IsNumber(peek) ||
			peek == '_' ||
			unicode.IsGraphic(peek) ||
			peek == '\u200D' || // Zero-Width Joiner
			peek == '\uFE0F' || // Variation Selector-16 (for emoji presentation)
			peek == '\uFE0E') {
			p.next()
		} else {
			stop := p.idx1
			token := string(p.Source[start:stop])

			tokenType, ok := p.keywords[token]
			if !ok {
				tokenType = Identifier
			}

			return Token{TokenType: tokenType, Token: token, Line: p.line}
		}
	}
}

func (p *Lexer) addToken(tokenType TokenType) {
	p.Tokens = append(p.Tokens, Token{
		TokenType: tokenType,
		Token:     string(p.Source[p.idx0:p.idx1]),
		Line:      p.line,
	})
}

func (p *Lexer) addTokenWithLiteral(tokenType TokenType, literal string) {
	p.Tokens = append(p.Tokens, Token{
		TokenType: tokenType,
		Token:     literal,
		Line:      p.line,
	})
}

func (p *Lexer) appendToken(token Token) {
	p.Tokens = append(p.Tokens, token)
}

func (p *Lexer) Tokenize() {
	p.line = 1
	p.column = 0

	c := p.chr0()

	if c != 0 {
		for {
			switch c {
			case '(':
				p.addToken(LeftParen)
			case ')':
				p.addToken(RightParen)
			case '{':
				p.addToken(LeftBrace)
			case '}':
				p.addToken(RightBrace)
			case '*':
				p.addToken(Star)
			case '.':
				p.addToken(Dot)
			case ',':
				p.addToken(Comma)
			case '+':
				p.addToken(Plus)
			case '-':
				p.addToken(Minus)
			case ';':
				p.addToken(Semicolon)

			case '!':
				if p.match('=') {
					p.addTokenWithLiteral(BangEqual, "!=")
					p.next()
				} else {
					p.addToken(Bang)
				}
			case '=':
				if p.match('=') {
					p.addTokenWithLiteral(EqualEqual, "==")
					p.next()
				} else {
					p.addToken(Equal)
				}
			case '>':
				if p.match('=') {
					p.addTokenWithLiteral(GreaterEqual, ">=")
					p.next()
				} else {
					p.addToken(Greater)
				}
			case '<':
				if p.match('=') {
					p.addTokenWithLiteral(LessEqual, "<=")
					p.next()
				} else {
					p.addToken(Less)
				}

			case '/':
				if p.match('/') {
					for {
						peek := p.peek()
						if peek == 0 || peek == '\n' {
							break
						} else {
							p.next()
						}
					}
				} else {
					p.addToken(Slash)
				}

			case '"':
				token, err := p.lexString()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v.\n", err)
				} else {
					p.appendToken(token)
				}

			default:
				if unicode.IsSpace(c) {
					// Nothing to do
				} else if unicode.IsDigit(c) {
					token, err := p.lexNumber()
					if err != nil {
						fmt.Fprintf(os.Stderr, "%v.\n", err)
					} else {
						p.appendToken(token)
					}
				} else if unicode.IsLetter(c) || c == '_' {
					token := p.lexIdentifer()
					p.appendToken(token)
				} else if SupportUnicode && !(c == '@' || c == '#' || c == '$' || c == '%' || c == '^' || c == '&' || c == '*') &&
					(unicode.IsLetter(c) || c == '_' || unicode.IsGraphic(c) ||
						c == '\u200D' || c == '\uFE0F' || c == '\uFE0E') {
					token := p.lexIdentifer()
					p.appendToken(token)
				} else {
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", p.line, string(c))
					p.HasLexicalErrors = true
				}
			}

			c = p.next()
			if c == 0 {
				break
			}
		}
	}
	p.addTokenWithLiteral(EOF, "")
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
	values := []string{t.TokenType.String(), t.Token, tokenData}
	return strings.Join(values, " ")
}

func (t TokenType) IsKeyword() bool {
	return t >= And && t <= While
}

func (t TokenType) IsLiteral() bool {
	return t == String || t == Number || t == Identifier
}

func (t TokenType) IsOperator() bool {
	return (t >= Star && t <= Slash) || (t >= Bang && t <= GreaterEqual)
}

func (t TokenType) String() string {
	switch t {
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
	case Number:
		return "NUMBER"
	case Identifier:
		return "IDENTIFIER"
	case And:
		return "AND"
	case Class:
		return "CLASS"
	case Else:
		return "ELSE"
	case False:
		return "FALSE"
	case For:
		return "FOR"
	case Fun:
		return "FUN"
	case If:
		return "IF"
	case Nil:
		return "NIL"
	case Or:
		return "OR"
	case Print:
		return "PRINT"
	case Return:
		return "RETURN"
	case Super:
		return "SUPER"
	case This:
		return "THIS"
	case True:
		return "TRUE"
	case Var:
		return "VAR"
	case While:
		return "WHILE"
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

	switch command {
	case "tokenize":
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		lexer := NewLexer(string(fileContents))
		lexer.Tokenize()
		lexer.PrintTokens()

		if lexer.HasLexicalErrors {
			os.Exit(LexicalError)
		}
	case "parse":

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
