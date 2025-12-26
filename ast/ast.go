// Package ast defines the Abstract Syntax Tree for Nulang.
package ast

import (
	"bytes"
	"strings"

	"github.com/nulang/nulang/token"
)

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents an identifier
type Identifier struct {
	Token  token.Token
	Value  string
	IsRest bool // true if this is a rest parameter (...args)
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	if i.IsRest {
		return "..." + i.Value
	}
	return i.Value
}

// LetStatement represents a let statement
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ConstStatement represents a const statement
type ConstStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")
	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// VarStatement represents a var statement
type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block statement { ... }
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// NumberLiteral represents a number literal
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// TemplateLiteral represents a template literal (backtick string with ${} expressions)
type TemplateLiteral struct {
	Token       token.Token
	Parts       []string     // Static string parts
	Expressions []Expression // Expressions between ${}
}

func (tl *TemplateLiteral) expressionNode()      {}
func (tl *TemplateLiteral) TokenLiteral() string { return tl.Token.Literal }
func (tl *TemplateLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("`")
	for i, part := range tl.Parts {
		out.WriteString(part)
		if i < len(tl.Expressions) {
			out.WriteString("${")
			out.WriteString(tl.Expressions[i].String())
			out.WriteString("}")
		}
	}
	out.WriteString("`")
	return out.String()
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// NullLiteral represents null
type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NullLiteral) String() string       { return "null" }

// UndefinedLiteral represents undefined
type UndefinedLiteral struct {
	Token token.Token
}

func (ul *UndefinedLiteral) expressionNode()      {}
func (ul *UndefinedLiteral) TokenLiteral() string { return ul.Token.Literal }
func (ul *UndefinedLiteral) String() string       { return "undefined" }

// PrefixExpression represents a prefix expression (!x, -x, ++x, --x)
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents an infix expression (x + y)
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression represents a postfix expression (x++, x--)
type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// IfExpression represents an if expression
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// ConditionalExpression represents a ternary expression (a ? b : c)
type ConditionalExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ce *ConditionalExpression) expressionNode()      {}
func (ce *ConditionalExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ce.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(ce.Consequence.String())
	out.WriteString(" : ")
	out.WriteString(ce.Alternative.String())
	out.WriteString(")")
	return out.String()
}

// FunctionLiteral represents a function literal
type FunctionLiteral struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	IsArrow    bool
	IsAsync    bool
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	if fl.IsAsync {
		out.WriteString("async ")
	}
	if fl.IsArrow {
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(") => ")
	} else {
		out.WriteString("function")
		if fl.Name != nil {
			out.WriteString(" " + fl.Name.String())
		}
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(") ")
	}
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents a function call
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral represents an array literal [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression represents array/object indexing
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// ObjectLiteral represents an object literal {a: 1, b: 2}
type ObjectLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) TokenLiteral() string { return ol.Token.Literal }
func (ol *ObjectLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range ol.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// MemberExpression represents property access (obj.prop or obj?.prop)
type MemberExpression struct {
	Token    token.Token
	Object   Expression
	Property Expression
	Optional bool // for optional chaining ?.
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	var out bytes.Buffer
	out.WriteString(me.Object.String())
	if me.Optional {
		out.WriteString("?.")
	} else {
		out.WriteString(".")
	}
	out.WriteString(me.Property.String())
	return out.String()
}

// AssignmentExpression represents an assignment (x = 5)
type AssignmentExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Left.String())
	out.WriteString(" " + ae.Operator + " ")
	out.WriteString(ae.Right.String())
	return out.String()
}

// ForStatement represents a for loop
type ForStatement struct {
	Token       token.Token
	Init        Statement
	Condition   Expression
	Update      Expression
	Body        *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for (")
	if fs.Init != nil {
		out.WriteString(fs.Init.String())
	}
	out.WriteString("; ")
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString("; ")
	if fs.Update != nil {
		out.WriteString(fs.Update.String())
	}
	out.WriteString(") ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// WhileStatement represents a while loop
type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while (")
	out.WriteString(ws.Condition.String())
	out.WriteString(") ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// BreakStatement represents break
type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return "break;" }

// ContinueStatement represents continue
type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return "continue;" }

// ThrowStatement represents throw
type ThrowStatement struct {
	Token token.Token
	Value Expression
}

func (ts *ThrowStatement) statementNode()       {}
func (ts *ThrowStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *ThrowStatement) String() string {
	return "throw " + ts.Value.String() + ";"
}

// TryStatement represents try/catch/finally
type TryStatement struct {
	Token        token.Token
	Block        *BlockStatement
	CatchParam   *Identifier
	CatchBlock   *BlockStatement
	FinallyBlock *BlockStatement
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	var out bytes.Buffer
	out.WriteString("try ")
	out.WriteString(ts.Block.String())
	if ts.CatchBlock != nil {
		out.WriteString(" catch")
		if ts.CatchParam != nil {
			out.WriteString(" (" + ts.CatchParam.String() + ")")
		}
		out.WriteString(" ")
		out.WriteString(ts.CatchBlock.String())
	}
	if ts.FinallyBlock != nil {
		out.WriteString(" finally ")
		out.WriteString(ts.FinallyBlock.String())
	}
	return out.String()
}

// NewExpression represents new ClassName()
type NewExpression struct {
	Token     token.Token
	Class     Expression
	Arguments []Expression
}

func (ne *NewExpression) expressionNode()      {}
func (ne *NewExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NewExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ne.Arguments {
		args = append(args, a.String())
	}
	out.WriteString("new ")
	out.WriteString(ne.Class.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ThisExpression represents 'this' keyword
type ThisExpression struct {
	Token token.Token
}

func (te *ThisExpression) expressionNode()      {}
func (te *ThisExpression) TokenLiteral() string { return te.Token.Literal }
func (te *ThisExpression) String() string       { return "this" }

// TypeofExpression represents typeof x
type TypeofExpression struct {
	Token token.Token
	Value Expression
}

func (te *TypeofExpression) expressionNode()      {}
func (te *TypeofExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TypeofExpression) String() string {
	return "typeof " + te.Value.String()
}

// AwaitExpression represents await x
type AwaitExpression struct {
	Token token.Token
	Value Expression
}

func (ae *AwaitExpression) expressionNode()      {}
func (ae *AwaitExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AwaitExpression) String() string {
	return "await " + ae.Value.String()
}

// SpreadExpression represents ...x
type SpreadExpression struct {
	Token token.Token
	Value Expression
}

func (se *SpreadExpression) expressionNode()      {}
func (se *SpreadExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SpreadExpression) String() string {
	return "..." + se.Value.String()
}

// ImportName represents a named import with optional alias
type ImportName struct {
	Token    token.Token
	Name     *Identifier // Original name from module
	Alias    *Identifier // Local name (if different from Name)
}

func (in *ImportName) String() string {
	if in.Alias != nil && in.Alias.Value != in.Name.Value {
		return in.Name.Value + " as " + in.Alias.Value
	}
	return in.Name.Value
}

// ImportStatement represents import statements
type ImportStatement struct {
	Token       token.Token
	Default     *Identifier       // import x from "..."
	Named       []*ImportName     // import { a, b } or { a as b } from "..."
	NamespaceAs *Identifier       // import * as x from "..."
	Source      *StringLiteral
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("import ")
	if is.Default != nil {
		out.WriteString(is.Default.String())
	}
	if len(is.Named) > 0 {
		names := []string{}
		for _, n := range is.Named {
			names = append(names, n.String())
		}
		out.WriteString("{ " + strings.Join(names, ", ") + " }")
	}
	if is.NamespaceAs != nil {
		out.WriteString("* as " + is.NamespaceAs.String())
	}
	out.WriteString(" from ")
	out.WriteString(is.Source.String())
	return out.String()
}

// ExportStatement represents export statements
type ExportStatement struct {
	Token     token.Token
	Default   Expression
	Named     []Statement
	Source    *StringLiteral
}

func (es *ExportStatement) statementNode()       {}
func (es *ExportStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("export ")
	if es.Default != nil {
		out.WriteString("default ")
		out.WriteString(es.Default.String())
	}
	return out.String()
}

// ClassDeclaration represents a class definition
type ClassDeclaration struct {
	Token      token.Token
	Name       *Identifier
	SuperClass Expression
	Implements []*Identifier // interface names that this class implements
	Decorators []*Decorator  // decorators applied to the class
	Body       *ClassBody
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) expressionNode()      {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ClassDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString("class ")
	if cd.Name != nil {
		out.WriteString(cd.Name.String())
	}
	if cd.SuperClass != nil {
		out.WriteString(" extends ")
		out.WriteString(cd.SuperClass.String())
	}
	if len(cd.Implements) > 0 {
		out.WriteString(" implements ")
		for i, iface := range cd.Implements {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(iface.String())
		}
	}
	out.WriteString(" ")
	out.WriteString(cd.Body.String())
	return out.String()
}

// ClassBody represents the body of a class
type ClassBody struct {
	Token   token.Token
	Members []ClassMember
}

func (cb *ClassBody) TokenLiteral() string { return cb.Token.Literal }
func (cb *ClassBody) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, member := range cb.Members {
		out.WriteString("  " + member.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// ClassMember represents a class member (method or property)
type ClassMember struct {
	Token      token.Token
	Name       *Identifier
	Value      Expression
	IsStatic   bool
	IsGetter   bool
	IsSetter   bool
	IsPrivate  bool
	Decorators []*Decorator // decorators applied to this member
}

func (cm *ClassMember) TokenLiteral() string { return cm.Token.Literal }
func (cm *ClassMember) String() string {
	var out bytes.Buffer
	if cm.IsStatic {
		out.WriteString("static ")
	}
	if cm.IsGetter {
		out.WriteString("get ")
	}
	if cm.IsSetter {
		out.WriteString("set ")
	}
	if cm.Name != nil {
		out.WriteString(cm.Name.String())
	}
	if cm.Value != nil {
		out.WriteString(" = ")
		out.WriteString(cm.Value.String())
	}
	return out.String()
}

// SuperExpression represents 'super' keyword
type SuperExpression struct {
	Token token.Token
}

func (se *SuperExpression) expressionNode()      {}
func (se *SuperExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SuperExpression) String() string       { return "super" }

// TypeAnnotation represents a type annotation (TypeScript-like)
type TypeAnnotation struct {
	Token    token.Token
	TypeName string
	IsArray  bool
	Generics []*TypeAnnotation
}

func (ta *TypeAnnotation) expressionNode()      {}
func (ta *TypeAnnotation) TokenLiteral() string { return ta.Token.Literal }
func (ta *TypeAnnotation) String() string {
	var out bytes.Buffer
	out.WriteString(ta.TypeName)
	if len(ta.Generics) > 0 {
		out.WriteString("<")
		for i, g := range ta.Generics {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(g.String())
		}
		out.WriteString(">")
	}
	if ta.IsArray {
		out.WriteString("[]")
	}
	return out.String()
}

// InterfaceDeclaration represents an interface definition
type InterfaceDeclaration struct {
	Token   token.Token
	Name    *Identifier
	Extends []*Identifier
	Body    []InterfaceMember
}

func (id *InterfaceDeclaration) statementNode()       {}
func (id *InterfaceDeclaration) TokenLiteral() string { return id.Token.Literal }
func (id *InterfaceDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString("interface ")
	out.WriteString(id.Name.String())
	if len(id.Extends) > 0 {
		out.WriteString(" extends ")
		for i, e := range id.Extends {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(e.String())
		}
	}
	out.WriteString(" { ... }")
	return out.String()
}

// InterfaceMember represents a member of an interface
type InterfaceMember struct {
	Token      token.Token
	Name       *Identifier
	Type       *TypeAnnotation
	IsOptional bool
	IsMethod   bool
	Parameters []*Identifier
}

func (im *InterfaceMember) TokenLiteral() string { return im.Token.Literal }
func (im *InterfaceMember) String() string {
	var out bytes.Buffer
	out.WriteString(im.Name.String())
	if im.IsOptional {
		out.WriteString("?")
	}
	if im.Type != nil {
		out.WriteString(": ")
		out.WriteString(im.Type.String())
	}
	return out.String()
}

// TypeAliasDeclaration represents a type alias
type TypeAliasDeclaration struct {
	Token token.Token
	Name  *Identifier
	Type  *TypeAnnotation
}

func (ta *TypeAliasDeclaration) statementNode()       {}
func (ta *TypeAliasDeclaration) TokenLiteral() string { return ta.Token.Literal }
func (ta *TypeAliasDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString("type ")
	out.WriteString(ta.Name.String())
	out.WriteString(" = ")
	if ta.Type != nil {
		out.WriteString(ta.Type.String())
	}
	return out.String()
}

// DeclareStatement represents a declare statement (TypeScript-like)
type DeclareStatement struct {
	Token token.Token
	Value Statement
}

func (ds *DeclareStatement) statementNode()       {}
func (ds *DeclareStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeclareStatement) String() string {
	var out bytes.Buffer
	out.WriteString("declare ")
	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}
	return out.String()
}

// Decorator represents a decorator @name or @name(args)
type Decorator struct {
	Token     token.Token
	Name      *Identifier
	Arguments []Expression
}

func (d *Decorator) expressionNode()      {}
func (d *Decorator) TokenLiteral() string { return d.Token.Literal }
func (d *Decorator) String() string {
	var out bytes.Buffer
	out.WriteString("@")
	out.WriteString(d.Name.String())
	if len(d.Arguments) > 0 {
		args := []string{}
		for _, a := range d.Arguments {
			args = append(args, a.String())
		}
		out.WriteString("(")
		out.WriteString(strings.Join(args, ", "))
		out.WriteString(")")
	}
	return out.String()
}
