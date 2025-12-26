// Package lexer implements the lexical analyzer for Nulang.
package lexer

import (
	"github.com/nulang/nulang/token"
)

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekCharN returns the character n positions ahead
func (l *Lexer) peekCharN(n int) byte {
	pos := l.readPosition + n - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	l.skipComments()
	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.EQ3, Literal: "===", Line: l.line, Column: l.column}
			} else {
				tok = token.Token{Type: token.EQ, Literal: "==", Line: l.line, Column: l.column}
			}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: "=>", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: "++", Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: "+=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.PLUS, l.ch, l.line, l.column)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: "--", Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: "-=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.MINUS, l.ch, l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.NOT_EQ3, Literal: "!==", Line: l.line, Column: l.column}
			} else {
				tok = token.Token{Type: token.NOT_EQ, Literal: "!=", Line: l.line, Column: l.column}
			}
		} else {
			tok = newToken(token.BANG, l.ch, l.line, l.column)
		}
	case '/':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.SLASH_ASSIGN, Literal: "/=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.SLASH, l.ch, l.line, l.column)
		}
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			tok = token.Token{Type: token.POWER, Literal: "**", Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.ASTERISK_ASSIGN, Literal: "*=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ASTERISK, l.ch, l.line, l.column)
		}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PERCENT_ASSIGN, Literal: "%=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.PERCENT, l.ch, l.line, l.column)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: "<=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: ">=", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.GT, l.ch, l.line, l.column)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: "&&", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: "||", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			tok = token.Token{Type: token.NULLISH, Literal: "??", Line: l.line, Column: l.column}
		} else if l.peekChar() == '.' {
			l.readChar()
			tok = token.Token{Type: token.OPTIONAL, Literal: "?.", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.QUESTION, l.ch, l.line, l.column)
		}
	case '.':
		if l.peekChar() == '.' && l.peekCharN(2) == '.' {
			l.readChar()
			l.readChar()
			tok = token.Token{Type: token.SPREAD, Literal: "...", Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.DOT, l.ch, l.line, l.column)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line, l.column)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.column)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.column)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.column)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
	case '"', '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
		tok.Line = l.line
		tok.Column = l.column
		l.readChar() // consume closing quote
		return tok
	case '`':
		tok.Type = token.TEMPLATE_STRING
		tok.Literal = l.readTemplateString()
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	
	// Handle integer part
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Handle decimal part
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	
	// Handle exponent part (e.g., 1e10, 1E-5)
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	
	// Handle BigInt suffix
	if l.ch == 'n' {
		l.readChar()
	}
	
	return l.input[position:l.position]
}

func (l *Lexer) readString(quote byte) string {
	var result []byte
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			case '\'':
				result = append(result, '\'')
			case '0':
				result = append(result, 0)
			default:
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
	}
	return string(result)
}

// readTemplateString reads a template literal (backtick string) with ${} expressions
// It returns the raw content including ${...} placeholders for later parsing
func (l *Lexer) readTemplateString() string {
	var result []byte
	l.readChar() // consume opening backtick
	for {
		if l.ch == '`' || l.ch == 0 {
			l.readChar() // consume closing backtick
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '`':
				result = append(result, '`')
			case '$':
				result = append(result, '$')
			case '0':
				result = append(result, 0)
			default:
				result = append(result, l.ch)
			}
			l.readChar()
		} else if l.ch == '$' && l.peekChar() == '{' {
			// Keep the ${...} as is for later parsing
			result = append(result, '$', '{')
			l.readChar() // consume $
			l.readChar() // consume {
			// Read until matching }
			depth := 1
			for depth > 0 && l.ch != 0 {
				if l.ch == '{' {
					depth++
				} else if l.ch == '}' {
					depth--
					if depth == 0 {
						result = append(result, '}')
						l.readChar()
						break
					}
				}
				if l.ch == '"' || l.ch == '\'' {
					// Handle strings inside expressions
					quote := l.ch
					result = append(result, l.ch)
					l.readChar()
					for l.ch != quote && l.ch != 0 {
						if l.ch == '\\' {
							result = append(result, l.ch)
							l.readChar()
						}
						result = append(result, l.ch)
						l.readChar()
					}
					result = append(result, l.ch)
					l.readChar()
				} else {
					result = append(result, l.ch)
					l.readChar()
				}
			}
		} else {
			result = append(result, l.ch)
			l.readChar()
		}
	}
	return string(result)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	// Single line comment
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
		l.skipComments()
	}
	
	// Multi-line comment
	if l.ch == '/' && l.peekChar() == '*' {
		l.readChar() // consume /
		l.readChar() // consume *
		for {
			if l.ch == '*' && l.peekChar() == '/' {
				l.readChar() // consume *
				l.readChar() // consume /
				break
			}
			if l.ch == 0 {
				break
			}
			l.readChar()
		}
		l.skipWhitespace()
		l.skipComments()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
