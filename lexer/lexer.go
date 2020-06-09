package lexer

import (
	"github.com/threeaccents/digolang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			tok.Literal = string(ch) + string(l.char)
			tok.Type = token.EQ
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '!':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			tok.Literal = string(ch) + string(l.char)
			tok.Type = token.NOT_EQ
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '[':
		tok = newToken(token.LBRACKET, l.char)
	case ']':
		tok = newToken(token.RBRACKET, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '.':
		tok = newToken(token.PERIOD, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIndentifier(tok.Literal)
			return tok
		}
		if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}

		tok = newToken(token.ILLEGAL, l.char)
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	defer func() {
		l.position = l.readPosition
		l.readPosition += 1
	}()

	if l.readPosition >= len(l.input) {
		l.char = 0
		return
	}

	l.char = l.input[l.readPosition]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) eatWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	l.readChar()

	position := l.position

	for l.char != '"' && l.char != 0 {
		l.readChar()
	}

	res := l.input[position:l.position]

	l.readChar()

	return res
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
