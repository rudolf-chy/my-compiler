package main

import (
	"fmt"
	"strconv"
)

type SimpleCalculator struct {
}

func (s *SimpleCalculator) prog(tokens TokenReader) *ASTNoder {
	noder := NewASTNoder(ASTNodeType_Program, "Calculator")
	child := s.intDeclare(tokens)
	if child != nil {
		noder.AddChild(*child)
	}
	return &noder
}

func (s *SimpleCalculator) intDeclare(reader TokenReader) *ASTNoder {
	var node ASTNoder
	token := reader.Peek()
	if token != nil && token.Type == TokenType_Int {
		reader.Read()
		if reader.Peek().Type == TokenType_Id {
			token := reader.Read()
			node = NewASTNoder(ASTNodeType_IntDeclaration, token.Text)
			token = reader.Peek()
			if token != nil && token.Type == TokenType_Assignment {
				reader.Read()
				child := s.addtive(reader)
				if child == nil {
					panic("invalide variable initialization, expecting an expression")
				} else {
					node.AddChild(*child)
				}
			}
		} else {
			panic("variable name expected")
		}
		// if node != nil {
		// 	token := reader.Peek()
		// 	if token != nil && token.Type == TokenType_SemiColon {
		// 		reader.Read()
		// 	} else {
		// 		panic("invalid statement, expecting semicolon")
		// 	}
		// }
	}
	return &node
}

func (s *SimpleCalculator) addtive(reader TokenReader) *ASTNoder {
	child1 := s.multiplicative(reader)
	var node ASTNoder
	if child1 != nil {
		for {
			token := reader.Peek()
			if token != nil && token.Type == TokenType_Plus {
				token := reader.Read()
				child2 := s.multiplicative(reader)
				if child2 != nil {
					node = NewASTNoder(ASTNodeType_AddtiveExp, token.Text)
					node.AddChild(*child1)
					node.AddChild(*child2)
					*child1 = node
				} else {
					panic("invalid additive expression, expecting the right part.")
				}
			} else {
				break
			}
		}

	}
	return child1
}

func (s *SimpleCalculator) multiplicative(reader TokenReader) *ASTNoder {
	child1 := s.primary(reader)
	var node ASTNoder
	token := reader.Peek()
	if child1 != nil && token != nil {
		if token.Type == TokenType_Star {
			token := reader.Read()
			child2 := s.primary(reader)
			if child2 != nil {
				node = NewASTNoder(ASTNodeType_Multiplicative, token.Text)
				node.AddChild(*child1)
				node.AddChild(*child2)
			} else {
				panic("invalid multiplicative expression, expecting the right part.")
			}
		}
	}
	if node != nil {
		return &node
	}
	return child1
}

func (s *SimpleCalculator) primary(reader TokenReader) *ASTNoder {
	var node ASTNoder
	token := reader.Peek()
	if token != nil {
		switch token.Type {
		case TokenType_IntLiteral:
			reader.Read()
			node = NewASTNoder(ASTNodeType_IntLiteral, token.Text)
		case TokenType_Id:
			token := reader.Read()
			node = NewASTNoder(ASTNodeType_IntDeclaration, token.Text)
		case TokenType_Left_Paren:
			reader.Read()
			node := s.addtive(reader)
			if node != nil {
				token := reader.Peek()
				if token != nil && token.Type == TokenType_Right_Paren {
					reader.Read()
				} else {
					panic("expecting right parenthesis")
				}
			} else {
				panic("expecting an additive expression inside parenthesis")
			}
		}
	}
	return &node
}

func (s *SimpleCalculator) Evaluate(script string) int {
	node := s.Parse(script)
	DumpAST(*node, "  ")
	return s.evaluate(*node, "  ")
}

func (s *SimpleCalculator) Parse(code string) *ASTNoder {
	lexer := SimpleLexer{}
	tokens := lexer.tokenize(code)
	return s.prog(tokens)
}

func (s *SimpleCalculator) evaluate(node ASTNoder, indent string) int {
	result := 0
	fmt.Printf("%sCalculating:%s\n", indent, node.GetType())
	switch node.GetType() {
	case ASTNodeType_Program:
		for _, n := range node.GetChildren() {
			result = s.evaluate(n, indent)
		}
	case ASTNodeType_AddtiveExp:
		child1 := node.GetChildren()[0]
		value1 := s.evaluate(child1, indent+"\t")
		child2 := node.GetChildren()[1]
		value2 := s.evaluate(child2, indent+"\t")
		if node.GetText() == "+" {
			result = value1 + value2
		} else {
			result = value1 - value2
		}
	case ASTNodeType_IntLiteral:
		result, _ = strconv.Atoi(node.GetText())
	case ASTNodeType_Multiplicative:
		child1 := node.GetChildren()[0]
		value1 := s.evaluate(child1, indent+"\t")
		child2 := node.GetChildren()[1]
		value2 := s.evaluate(child2, indent+"\t")
		if node.GetText() == "*" {
			result = value1 * value2
		} else {
			result = value1 / value2
		}
	case ASTNodeType_IntDeclaration:
		for _, n := range node.GetChildren() {
			result = s.evaluate(n, indent)
		}
	}
	return result
}
