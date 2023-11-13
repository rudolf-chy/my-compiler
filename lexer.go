package main

import (
	"bytes"
	"fmt"
)

type DfaState int

const (
	DfaState_Initial DfaState = iota
	DfaState_Id
	DfaState_Int1
	DfaState_Int2
	DfaState_Int3
	DfaState_Assignment
	DfaState_SemiColon
	DfaState_Left_Paren
	DfaState_Right_Paren
	DfaState_GT
	DfaState_GE
	DfaState_LT
	DfaState_LE
	DfaState_Plus
	DfaState_Minus
	DfaState_Star
	DfaState_Slash
	DfaState_IntLiteral
)

type TokenType string

const (
	TokenType_Initial     = TokenType("Initial")
	TokenType_Id          = TokenType("Identifier")
	TokenType_GT          = TokenType("GT")
	TokenType_GE          = TokenType("GE")
	TokenType_LT          = TokenType("LT")
	TokenType_LE          = TokenType("LE")
	TokenType_IntLiteral  = TokenType("IntLiteral")
	TokenType_Int         = TokenType("Int")
	TokenType_Assignment  = TokenType("Assignment")
	TokenType_SemiColon   = TokenType("SemiColon")
	TokenType_Plus        = TokenType("Plus")
	TokenType_Minus       = TokenType("Minus")
	TokenType_Star        = TokenType("Star")
	TokenType_Slash       = TokenType("Slash")
	TokenType_Left_Paren  = TokenType("(")
	TokenType_Right_Paren = TokenType(")")
)

type TokenReader interface {
	Read() *Token
	Peek() *Token
	UnRead()
	GetPosition() int
	setPosition(position int)
}

func isAlpha(ch int32) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isDigit(ch int32) bool {
	return ch >= '0' && ch <= '9'
}

type Token struct {
	Text string
	Type TokenType
}

type SimpleLexer struct {
	tokens    []Token
	token     Token
	tokenText *bytes.Buffer
}

func NewSimpleLexer() SimpleLexer {
	return SimpleLexer{}
}

func (s *SimpleLexer) tokenize(script string) TokenReader {
	s.tokenText = new(bytes.Buffer)
	state := DfaState_Initial
	for _, ch := range script {
		switch state {
		case DfaState_Initial:
			state = s.initToken(ch)
		case DfaState_Id:
			if isAlpha(ch) || isDigit(ch) {
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaState_Int1:
			if ch == 'n' {
				state = DfaState_Int2
				s.tokenText.WriteRune(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				s.tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = s.initToken(ch)
			}
		case DfaState_Int2:
			if ch == 't' {
				state = DfaState_Int3
				s.tokenText.WriteRune(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				s.tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = s.initToken(ch)
			}
		case DfaState_Int3:
			if ch == ' ' {
				s.token.Type = TokenType_Int
				state = s.initToken(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				s.tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = s.initToken(ch)
			}

		case DfaState_GT:
			if ch == '=' {
				s.token.Type = TokenType_GE
				state = DfaState_GE
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaState_GE:
			state = s.initToken(ch)
		case DfaState_LT:
			if ch == '=' {
				s.token.Type = TokenType_LE
				state = DfaState_LE
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaState_LE:
			state = s.initToken(ch)
		case DfaState_Assignment:
			state = s.initToken(ch)
		case DfaState_Plus:
			state = s.initToken(ch)
		case DfaState_Minus:
			state = s.initToken(ch)
		case DfaState_Star:
			state = s.initToken(ch)
		case DfaState_Slash:
			state = s.initToken(ch)
		case DfaState_SemiColon:
			state = s.initToken(ch)
		case DfaState_Left_Paren:
			state = s.initToken(ch)
		case DfaState_Right_Paren:
			state = s.initToken(ch)
		case DfaState_IntLiteral:
			if isDigit(ch) {
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		}
	}
	return NewTokenReader(s.tokens)
}

func (s *SimpleLexer) initToken(ch rune) DfaState {
	if len(s.tokenText.Bytes()) > 0 {
		s.token.Text = s.tokenText.String()
		s.tokens = append(s.tokens, s.token)
	}
	s.tokenText = new(bytes.Buffer)
	s.token = Token{}
	newstate := DfaState_Initial
	switch {
	case isAlpha(ch):
		if ch == 'i' {
			newstate = DfaState_Int1
		} else {
			newstate = DfaState_Id
		}
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Id
	case isDigit(ch):
		newstate = DfaState_IntLiteral
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_IntLiteral
	case ch == '<':
		newstate = DfaState_LT
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_LT
	case ch == '>':
		newstate = DfaState_GT
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_GT
	case ch == '=':
		newstate = DfaState_Assignment
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Assignment
	case ch == '+':
		newstate = DfaState_Plus
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Plus
	case ch == '-':
		newstate = DfaState_Minus
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Minus
	case ch == '*':
		newstate = DfaState_Star
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Star
	case ch == '/':
		newstate = DfaState_Slash
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Slash
	case ch == ';':
		newstate = DfaState_SemiColon
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_SemiColon
	case ch == '(':
		newstate = DfaState_Left_Paren
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Left_Paren
	case ch == ')':
		newstate = DfaState_Right_Paren
		s.tokenText.WriteRune(ch)
		s.token.Type = TokenType_Right_Paren
	}
	return newstate
}

func (lexer *SimpleLexer) dump(reader TokenReader) {
	fmt.Println("text\ttype")
	var token *Token
	for {
		if token = reader.Read(); token == nil {
			break
		} else {
			fmt.Printf("%s\t\t%s\n", (*token).Text, (*token).Type)
		}
	}
}

type SimpleTokenReader struct {
	tokens   []Token
	position int
}

func (s *SimpleTokenReader) Read() *Token {
	if s.position < len(s.tokens) {
		p := s.position
		s.position++
		return &s.tokens[p]
	}
	return nil
}

func (s *SimpleTokenReader) Peek() *Token {
	if s.position < len(s.tokens) {
		p := s.position
		return &s.tokens[p]
	}
	return nil
}

func (s *SimpleTokenReader) UnRead() {
	if s.position > 0 {
		s.position--
	}
}

func (s *SimpleTokenReader) GetPosition() int {
	return s.position
}

func (s *SimpleTokenReader) setPosition(position int) {
	s.position = position
}

func NewTokenReader(tokens []Token) TokenReader {
	return &SimpleTokenReader{tokens: tokens}
}
