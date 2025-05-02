package ast

import (
	"context"
	"strconv"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Expression interface {
	Node
	isExpression()
}

type Expressions []Expression

func (e Expressions) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, e, ", ")
}

var NoArguments = Expressions{}

type BuiltinExpression string

const (
	ExpressionThis   BuiltinExpression = "this"
	ExpressionOuter  BuiltinExpression = "outer"
	ExpressionModule BuiltinExpression = "module"
	ExpressionNull   BuiltinExpression = "null"
	ExpressionTrue   BuiltinExpression = "true"
	ExpressionFalse  BuiltinExpression = "false"
)

func (e BuiltinExpression) isExpression() {}

func (e BuiltinExpression) Marshal(_ context.Context) ([]byte, error) {
	return []byte(e), nil
}

// IntExpression represents an integer literal.
type IntExpression int64

func (e IntExpression) isExpression() {}

func (e IntExpression) Marshal(_ context.Context) ([]byte, error) {
	return []byte(strconv.FormatInt(int64(e), 10)), nil
}

// FloatExpression represents a floating-point literal.
type FloatExpression float64

func (e FloatExpression) isExpression() {}

func (e FloatExpression) Marshal(_ context.Context) ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(e), 'f', -1, 64)), nil
}

// StringExpression represents a string literal.
type StringExpression string

func (e StringExpression) isExpression() {}

func (e StringExpression) Marshal(_ context.Context) ([]byte, error) {
	return []byte(strconv.Quote(string(e))), nil
}

type PrefixUnaryExpression struct {
	// Required.
	Operator PrefixUnaryOperand
	// Required.
	Operand Expression
}

func (e *PrefixUnaryExpression) isExpression() {}

func (e *PrefixUnaryExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	operand, err := e.Operand.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	b.WriteString(string(e.Operator))
	b.Write(operand)

	return b.Bytes(), nil
}

type PostfixUnaryExpression struct {
	// Required.
	Operator PostfixUnaryOperand
	// Required.
	Operand Expression
}

func (e *PostfixUnaryExpression) isExpression() {}

func (e *PostfixUnaryExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	operand, err := e.Operand.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	b.Write(operand)
	b.WriteString(string(e.Operator))

	return b.Bytes(), nil
}

type BinaryExpression struct {
	// Required.
	Operator BinaryOperator
	// Required.
	Left Expression
	// Required.
	Right Expression
}

func (e *BinaryExpression) isExpression() {}

func (e *BinaryExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	left, err := e.Left.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	right, err := e.Right.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	b.WriteWithSuffix(left, " ")
	b.WriteString(string(e.Operator))
	b.WriteWithPrefix(" ", right)

	return b.Bytes(), nil
}

type TypeExpression struct {
	// Required.
	Operator TypeOperator
	// Required.
	Expression Expression
	// Required.
	Type Type
}

func (e *TypeExpression) isExpression() {}

func (e *TypeExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	expr, err := e.Expression.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	typ, err := e.Type.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	b.WriteWithSuffix(expr, " ")
	b.WriteString(string(e.Operator))
	b.WriteWithPrefix(" ", typ)

	return b.Bytes(), nil
}

type MemberAccessExpression struct {
	// Required.
	Name Identifier
	// Optional.
	Arguments Expressions
}

func (e *MemberAccessExpression) isExpression() {}

func (e *MemberAccessExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	name, err := e.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(name)

	if e.Arguments != nil {
		args, err := e.Arguments.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteRune('(')
		b.Write(args)
		b.WriteRune(')')
	}

	return b.Bytes(), nil
}

type QualifiedMemberAccessExpression struct {
	// Required.
	Receiver Expression
	// Optional.
	Nullable bool
	// Required.
	Name Identifier
	// Optional.
	Arguments Expressions
}

func (e *QualifiedMemberAccessExpression) isExpression() {}

func (e *QualifiedMemberAccessExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	receiver, err := e.Receiver.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(receiver)
	if e.Nullable {
		b.WriteString(".?")
	} else {
		b.WriteRune('.')
	}

	name, err := e.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(name)

	if e.Arguments != nil {
		args, err := e.Arguments.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteRune('(')
		b.Write(args)
		b.WriteRune(')')
	}

	return b.Bytes(), nil
}

type SuperAccessExpression struct {
	// Required.
	Name Identifier
	// Optional.
	Arguments Expressions
}

func (e *SuperAccessExpression) isExpression() {}

func (e *SuperAccessExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	name, err := e.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix("super.", name)

	if e.Arguments != nil {
		args, err := e.Arguments.Marshal(ctx)
		if err != nil {
			return nil, err
		}

		b.WriteRune('(')
		b.Write(args)
		b.WriteRune(')')
	}

	return b.Bytes(), nil
}

type SubscriptExpression struct {
	// Required.
	Receiver Expression
	// Required.
	Subscript Expression
}

func (e *SubscriptExpression) isExpression() {}

func (e *SubscriptExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	receiver, err := e.Receiver.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(receiver)

	subscript, err := e.Subscript.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("[", subscript, "]")

	return b.Bytes(), nil
}

type SuperSubscriptExpression struct {
	// Required.
	Subscript Expression
}

func (e *SuperSubscriptExpression) isExpression() {}

func (e *SuperSubscriptExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	subscript, err := e.Subscript.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("super[", subscript, "]")

	return b.Bytes(), nil
}

type ParenthesizedExpression struct {
	// Required.
	Expression Expression
}

func (e *ParenthesizedExpression) isExpression()            {}
func (e *ParenthesizedExpression) isAmendParentExpression() {}

func (e *ParenthesizedExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	expr, err := e.Expression.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("(", expr, ")")

	return b.Bytes(), nil
}

type NewExpression struct {
	// Optional.
	Type Type
	// Required.
	Body *ObjectBody
}

func (e *NewExpression) isExpression()            {}
func (e *NewExpression) isAmendParentExpression() {}

func (e *NewExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	b.WriteString("new ")

	if e.Type != nil {
		typ, err := e.Type.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithSuffix(typ, " ")
	}

	body, err := e.Body.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(body)

	return b.Bytes(), nil
}

type AmendParentExpression interface {
	Node
	isAmendParentExpression()
}

type AmendExpression struct {
	// Required.
	Parent AmendParentExpression
	// Required.
	Body *ObjectBody
}

func (e *AmendExpression) isExpression()            {}
func (e *AmendExpression) isAmendParentExpression() {}

func (e *AmendExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	parent, err := e.Parent.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	body, err := e.Body.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	b.Write(parent)
	b.WriteRune(' ')
	b.Write(body)

	return b.Bytes(), nil
}

type IfExpression struct {
	// Required.
	Condition Expression
	// Required.
	Then Expression
	// Required.
	Else Expression
}

func (e *IfExpression) isExpression() {}

func (e *IfExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	cond, err := e.Condition.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("if (", cond, ") ")

	then, err := e.Then.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(then)

	elseExpr, err := e.Else.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" else ", elseExpr)

	return b.Bytes(), nil
}

type ImportExpression struct {
	// Required.
	Path string
	// Optional.
	Glob bool
}

func (e *ImportExpression) isExpression() {}

func (e *ImportExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	if e.Glob {
		b.WriteString("import*(")
	} else {
		b.WriteString("import(")
	}
	b.WriteString(strconv.Quote(e.Path))
	b.WriteRune(')')

	return b.Bytes(), nil
}

type LetExpression struct {
	// Required.
	Name *Parameter
	// Required.
	Value Expression
	// Required.
	Expression Expression
}

func (e *LetExpression) isExpression() {}

func (e *LetExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	name, err := e.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("let (", name, " = ")

	val, err := e.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(val, ") ")

	expr, err := e.Expression.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(expr)

	return b.Bytes(), nil
}

type ReadVariant string

const (
	ReadVariantNullable ReadVariant = "nullable"
	ReadVariantGlob     ReadVariant = "glob"
)

type ReadExpression struct {
	// Optional.
	Variant ReadVariant
	// Required.
	Value Expression
}

func (e *ReadExpression) isExpression() {}

func (e *ReadExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	switch e.Variant {
	case ReadVariantNullable:
		b.WriteString("read?(")
	case ReadVariantGlob:
		b.WriteString("read*(")
	default:
		b.WriteString("read(")
	}

	expr, err := e.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(expr)
	b.WriteRune(')')

	return b.Bytes(), nil
}

type ThrowExpression struct {
	// Required.
	Value Expression
}

func (e *ThrowExpression) isExpression() {}

func (e *ThrowExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	expr, err := e.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("throw(", expr, ")")

	return b.Bytes(), nil
}

type TraceExpression struct {
	// Required.
	Value Expression
}

func (e *TraceExpression) isExpression() {}

func (e *TraceExpression) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	expr, err := e.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("trace(", expr, ")")

	return b.Bytes(), nil
}
