package ast

type Expr interface {
	expr()
}

type ExprImpl struct {
}

func (x *ExprImpl) expr() {}

type NumberExpr struct {
	ExprImpl
	Lit string
}

type StringExpr struct {
	ExprImpl
	Lit string
}

type ArrayExpr struct {
	ExprImpl
	Exprs []Expr
}

type IdentExpr struct {
	ExprImpl
	Lit string
}

type UnaryMinusExpr struct {
	ExprImpl
	SubExpr Expr
}

type ParenExpr struct {
	ExprImpl
	SubExpr Expr
}

type BinOpExpr struct {
	ExprImpl
	Lhs      Expr
	Operator string
	Rhs      Expr
}

type CallExpr struct {
	ExprImpl
	Name     string
	SubExprs []Expr
}