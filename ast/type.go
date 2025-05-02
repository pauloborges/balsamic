package ast

import (
	"context"
	"strconv"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Type interface {
	Node
	isType()
}

type BuiltinType string

const (
	TypeUnknown BuiltinType = "unknown"
	TypeNothing BuiltinType = "nothing"
	TypeModule  BuiltinType = "module"
)

func (t BuiltinType) isType() {}

func (t BuiltinType) Marshal(_ context.Context) ([]byte, error) {
	return []byte(t), nil
}

type StringLiteralType string

func (t StringLiteralType) isType() {}

func (t StringLiteralType) Marshal(_ context.Context) ([]byte, error) {
	return []byte(strconv.Quote(string(t))), nil
}

type DeclaredType struct {
	Name           QualifiedIdentifier
	TypeParameters []Type
}

func (t *DeclaredType) isType() {}

func (t *DeclaredType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	name, err := t.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(name)

	if len(t.TypeParameters) > 0 {
		params, err := joinNodes(ctx, t.TypeParameters, ", ")
		if err != nil {
			return nil, err
		}

		b.WriteRune('<')
		b.Write(params)
		b.WriteByte('>')
	}

	return b.Bytes(), nil
}

type ParenthesizedType struct {
	// Required.
	Type Type
}

func (t *ParenthesizedType) isType() {}

func (t *ParenthesizedType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	typ, err := t.Type.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("(", typ, ")")

	return b.Bytes(), nil
}

type NullableType struct {
	// Required.
	Type Type
}

func (t *NullableType) isType() {}

func (t *NullableType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	typ, err := t.Type.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(typ, "?")

	return b.Bytes(), nil
}

type ConstrainedType struct {
	// Required.
	Type Type
	// Optional.
	Constraints Expressions
}

func (t *ConstrainedType) isType() {}

func (t *ConstrainedType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	typ, err := t.Type.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(typ)

	constraints, err := t.Constraints.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("(", constraints, ")")

	return b.Bytes(), nil
}

type UnionType struct {
	// Required.
	Members []Type
	// Optional.
	Default Type
}

func (t *UnionType) isType() {}

func (t *UnionType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	members, err := joinNodes(ctx, t.Members, " | ")
	if err != nil {
		return nil, err
	}
	b.Write(members)

	if t.Default != nil {
		dflt, err := t.Default.Marshal(ctx)
		if err != nil {
			return nil, err
		}

		if len(t.Members) > 0 {
			b.WriteString(" | *")
		} else {
			b.WriteRune('*')
		}
		b.Write(dflt)
	}

	return b.Bytes(), nil
}

type FunctionLiteralType struct {
	// Optional.
	Parameters []Type
	// Required.
	Result Type
}

func (t *FunctionLiteralType) isType() {}

func (t *FunctionLiteralType) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	params, err := joinNodes(ctx, t.Parameters, ", ")
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("(", params, ")")

	result, err := t.Result.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" -> ", result)

	return b.Bytes(), nil
}
