package main

import (
	"fmt"
	"strconv"
)

/**
 * 一个简单的语法解析器。
 * 能够解析简单的表达式、变量声明和初始化语句、赋值语句。
 * 它支持的语法规则为：
 *
 * programm -> intDeclare | expressionStatement | assignmentStatement
 * intDeclare -> 'int' Id ( = additive) ';'
 * expressionStatement -> addtive ';'
 * addtive -> multiplicative ( (+ | -) multiplicative)*
 * multiplicative -> primary ( (* | /) primary)*
 * primary -> IntLiteral | Id | (additive)
 */

type ASTNodeType string

const (
	ASTNodeType_Program        = ASTNodeType("program")
	ASTNodeType_IntLiteral     = ASTNodeType("IntLiteral")
	ASTNodeType_IntDeclaration = ASTNodeType("IntDeclaration")
	ASTNodeType_AddtiveExp     = ASTNodeType("AddtiveExp")
	ASTNodeType_Multiplicative = ASTNodeType("Multiplicative")
	ASTNodeType_Assignment     = ASTNodeType("Assignment")
	ASTNodeType_Identifier     = ASTNodeType("Identifier")
)

type ASTNoder interface {
	AddChild(child ASTNoder)
	GetText() string
	GetType() ASTNodeType
	GetChildren() []ASTNoder
	GetParent() ASTNoder
}

type SimpleASTNode struct {
	nodeType ASTNodeType
	text     string
	parent   ASTNoder
	children []ASTNoder
}

func NewASTNoder(nodeType ASTNodeType, text string) ASTNoder {
	return &SimpleASTNode{nodeType: nodeType, text: text}
}

func (s *SimpleASTNode) AddChild(child ASTNoder) {
	s.children = append(s.children, child)
}

func (s *SimpleASTNode) GetText() string {
	return s.text
}

func (s *SimpleASTNode) GetType() ASTNodeType {
	return s.nodeType
}

func (s *SimpleASTNode) GetChildren() []ASTNoder {
	return s.children
}

func (s *SimpleASTNode) GetParent() ASTNoder {
	return s.parent
}

type SimpleParser struct {
}

func DumpAST(node ASTNoder, indent string) {
	fmt.Printf("%s%s %s\n", indent, node.GetType(), node.GetText())
	for _, _node := range node.GetChildren() {
		DumpAST(_node, "\t"+indent)
	}
}

func (s *SimpleParser) Evaluate(script string) int {
	node := s.Parse(script)
	DumpAST(*node, "	")
	return s.evaluate(*node, "\t")
}

func (s *SimpleParser) Parse(code string) *ASTNoder {
	lexer := SimpleLexer{}
	tokens := lexer.tokenize(code)
	return s.prog(tokens)
}

func (s *SimpleParser) evaluate(node ASTNoder, indent string) int {
	result := 0
	fmt.Printf("%sCalculating:%s\n", indent, node.GetType())
	switch node.GetType() {
	case ASTNodeType_Program:
		for _, n := range node.GetChildren() {
			result += s.evaluate(n, indent)
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
	}
	return result
}

func (s *SimpleParser) prog(tokens TokenReader) *ASTNoder {
	noder := NewASTNoder(ASTNodeType_Program, "pwc")
	for {
		token := tokens.Peek()
		if token == nil {
			break
		}
		child := s.intDeclare(tokens)
		if child == nil {
			child = s.expressionStatement(tokens)
		}
		if child == nil {
			child = s.assignmentStatement(tokens)
		}
		if child != nil {
			noder.AddChild(*child)
		} else {
			panic("unknown statement")
		}
	}

	return &noder
}

func (s *SimpleParser) intDeclare(reader TokenReader) *ASTNoder {
	var node ASTNoder
	token := reader.Peek()
	if token != nil && token.Type == TokenType_Int {
		reader.Read()
		if reader.Peek().Type == TokenType_Id {
			token = reader.Read()
			node = NewASTNoder(ASTNodeType_IntDeclaration, token.Text)
			token = reader.Peek()
			if token != nil && token.Type == TokenType_Assignment {
				reader.Read()
				child := s.additive(reader)
				if child == nil {
					panic("invalide variable initialization, expecting an expression")
				} else {
					node.AddChild(*child)
				}
			}
		} else {
			panic("variable name expected")
		}
		if node != nil {
			token = reader.Peek()
			if token != nil && token.Type == TokenType_SemiColon {
				reader.Read()
			} else {
				panic("variable name expected")
			}
		}
	}
	if node != nil {
		return &node
	}
	return nil
}

func (s *SimpleParser) additive1(reader TokenReader) *ASTNoder {
	child1 := s.multiplicative(reader)
	var node ASTNoder
	token := reader.Peek()
	if child1 != nil && token != nil {
		if token.Type == TokenType_Plus {
			token := reader.Read()
			child2 := s.additive1(reader)
			if child2 != nil {
				node = NewASTNoder(ASTNodeType_AddtiveExp, token.Text)
				node.AddChild(*child1)
				node.AddChild(*child2)
			} else {
				panic("invalid additive expression, expecting the right part.")
			}
		}
	}
	if node != nil {
		return &node
	}
	return child1
}

func (s *SimpleParser) additive(reader TokenReader) *ASTNoder {
	child1 := s.multiplicative(reader)
	var node ASTNoder
	if child1 != nil {
		for {
			token := reader.Peek()
			if token != nil && token.Type == TokenType_Plus {
				token = reader.Read()
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

func (s *SimpleParser) expressionStatement(reader TokenReader) *ASTNoder {
	pos := reader.GetPosition()
	node := s.additive(reader)
	if node != nil {
		token := reader.Peek()
		if token != nil && token.Type == TokenType_SemiColon {
			reader.Read()
		} else {
			node = nil
			reader.setPosition(pos)
		}
	}
	return node
}

func (s *SimpleParser) assignmentStatement(reader TokenReader) *ASTNoder {
	pos := reader.GetPosition()
	var node ASTNoder
	token := reader.Peek()
	if token != nil && token.Type == TokenType_Id {
		token = reader.Read()
		node = NewASTNoder(ASTNodeType_Assignment, token.Text)
		token = reader.Peek()
		if token != nil && token.Type == TokenType_Assignment {
			reader.Read()
			child := s.additive(reader)
			if child != nil {
				node.AddChild(*child)
				token = reader.Peek()
				if token != nil && token.Type == TokenType_SemiColon {
					reader.Read()
				} else {
					panic("invalid statement, expecting semicolon")
				}
			} else {
				panic("invalide assignment statement, expecting an expression")
			}
		} else {
			reader.setPosition(pos)
			return nil
		}
	}
	return &node
}

func (s *SimpleParser) primary(reader TokenReader) *ASTNoder {
	var node ASTNoder
	token := reader.Peek()
	if token != nil {
		switch token.Type {
		case TokenType_IntLiteral:
			reader.Read()
			node = NewASTNoder(ASTNodeType_IntLiteral, token.Text)
		case TokenType_Id:
			reader.Read()
			node = NewASTNoder(ASTNodeType_Identifier, token.Text)
		case TokenType_Left_Paren:
			reader.Read()
			node := s.additive(reader)
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

func (s *SimpleParser) multiplicative(reader TokenReader) *ASTNoder {
	child1 := s.primary(reader)
	var node ASTNoder
	token := reader.Peek()
	if child1 != nil && token != nil {
		if token.Type == TokenType_Star {
			token = reader.Read()
			child2 := s.primary(reader)
			if child2 != nil {
				node = NewASTNoder(ASTNodeType_Multiplicative, token.Text)
				node.AddChild(*child1)
				node.AddChild(*child2)
				child1 = &node
			} else {
				panic("invalid multiplicative expression, expecting the right part.")
			}
		}
	}
	return child1
}
