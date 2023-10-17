package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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

// type SimpleLexer struct {
// 	tokens    []Token
// 	token     Token
// 	tokenText *bytes.Buffer
// }

var tokens []Token
var token Token
var tokenText *bytes.Buffer

func main() {
	// file, err := os.OpenFile(os.Args[1], os.O_RDWR, 0666)
	inputReader := bufio.NewReader(os.Stdin)
	script, _ := inputReader.ReadString('\n')
	tokenText = new(bytes.Buffer)

	state := DfaState_Initial
	for _, ch := range script {
		switch state {
		case DfaState_Initial:
			state = initToken(ch)
		case DfaState_Id:
			if isAlpha(ch) || isDigit(ch) {
				tokenText.WriteRune(ch)
			} else {
				state = initToken(ch)
			}
		case DfaState_Int1:
			if ch == 'n' {
				state = DfaState_Int2
				tokenText.WriteRune(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = initToken(ch)
			}
		case DfaState_Int2:
			if ch == 't' {
				state = DfaState_Int3
				tokenText.WriteRune(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = initToken(ch)
			}
		case DfaState_Int3:
			if ch == ' ' {
				token.Type = TokenType_Int
				state = initToken(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				tokenText.WriteRune(ch)
				state = DfaState_Id
			} else {
				state = initToken(ch)
			}

		case DfaState_GT:
			if ch == '=' {
				token.Type = TokenType_GE
				state = DfaState_GE
				tokenText.WriteRune(ch)
			} else {
				state = initToken(ch)
			}
		case DfaState_GE:
			state = initToken(ch)
		case DfaState_LT:
			if ch == '=' {
				token.Type = TokenType_LE
				state = DfaState_LE
				tokenText.WriteRune(ch)
			} else {
				state = initToken(ch)
			}
		case DfaState_LE:
			state = initToken(ch)
		case DfaState_Assignment:
			state = initToken(ch)
		case DfaState_Plus:
			state = initToken(ch)
		case DfaState_Minus:
			state = initToken(ch)
		case DfaState_Star:
			state = initToken(ch)
		case DfaState_Slash:
			state = initToken(ch)
		case DfaState_SemiColon:
			state = initToken(ch)
		case DfaState_Left_Paren:
			state = initToken(ch)
		case DfaState_Right_Paren:
			state = initToken(ch)
		case DfaState_IntLiteral:
			if isDigit(ch) {
				tokenText.WriteRune(ch)
			} else {
				state = initToken(ch)
			}
		}
	}

	fmt.Println("tokens: ", tokens)
}

func initToken(ch rune) DfaState {
	if len(tokenText.Bytes()) > 0 {
		token.Text = tokenText.String()
		tokens = append(tokens, token)
	}
	tokenText = new(bytes.Buffer)
	token = Token{}
	newstate := DfaState_Initial
	switch {
	case isAlpha(ch):
		if ch == 'i' {
			newstate = DfaState_Int1
		} else {
			newstate = DfaState_Id
		}
		tokenText.WriteRune(ch)
		token.Type = TokenType_Id
	case isDigit(ch):
		newstate = DfaState_IntLiteral
		tokenText.WriteRune(ch)
		token.Type = TokenType_IntLiteral
	case ch == '<':
		newstate = DfaState_LT
		tokenText.WriteRune(ch)
		token.Type = TokenType_LT
	case ch == '>':
		newstate = DfaState_GT
		tokenText.WriteRune(ch)
		token.Type = TokenType_GT
	case ch == '=':
		newstate = DfaState_Assignment
		tokenText.WriteRune(ch)
		token.Type = TokenType_Assignment
	case ch == '+':
		newstate = DfaState_Plus
		tokenText.WriteRune(ch)
		token.Type = TokenType_Plus
	case ch == '-':
		newstate = DfaState_Minus
		tokenText.WriteRune(ch)
		token.Type = TokenType_Minus
	case ch == '*':
		newstate = DfaState_Star
		tokenText.WriteRune(ch)
		token.Type = TokenType_Star
	case ch == '/':
		newstate = DfaState_Slash
		tokenText.WriteRune(ch)
		token.Type = TokenType_Slash
	case ch == ';':
		newstate = DfaState_SemiColon
		tokenText.WriteRune(ch)
		token.Type = TokenType_SemiColon
	case ch == '(':
		newstate = DfaState_Left_Paren
		tokenText.WriteRune(ch)
		token.Type = TokenType_Left_Paren
	case ch == ')':
		newstate = DfaState_Right_Paren
		tokenText.WriteRune(ch)
		token.Type = TokenType_Right_Paren
	}
	return newstate
}
