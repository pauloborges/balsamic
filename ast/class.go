package ast

import (
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Class struct {
	// Optional.
	Docs Docs
	// Optional.
	Annotations Annotations
	// Optional.
	Modifiers Modifiers
	// Required.
	Name Identifier
	// Optional.
	TypeParameters TypeParameters
	// Optional.
	ParentName QualifiedIdentifier
	// Optional.
	ParentTypeParameters TypeParameters
	// Optional.
	Members []ClassMember
}

func (c *Class) isModuleMember() {}

func (c *Class) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	docs, err := c.Docs.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(docs, newlineWithIndentation(ctx))

	annotations, err := c.Annotations.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(annotations, newlineWithIndentation(ctx))

	modifiers, err := c.Modifiers.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(modifiers, " ")

	name, err := c.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix("class ", name)

	typeParams, err := c.TypeParameters.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(typeParams)

	parentName, err := c.ParentName.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" extends ", parentName)

	parentTypeParams, err := c.ParentTypeParameters.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(parentTypeParams)

	membersCtx := raiseIndentLevel(ctx)
	membersSep := "\n" + newlineWithIndentation(membersCtx)
	members, err := joinNodes(membersCtx, c.Members, membersSep)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix(
		" {"+newlineWithIndentation(membersCtx),
		members,
		newlineWithIndentation(ctx)+"}",
	)

	return b.Bytes(), nil
}

type ClassMember interface {
	Node
	isClassMember()
}

// ClassProperty represents a property declaration in a class or module.
//
// It can either be a property with a type and/or an expression, or a
// property with a body (ammending).
//
// Setting both type/expression and body is not allowed.
type ClassProperty struct {
	Docs        Docs
	Annotations Annotations
	Modifiers   Modifiers
	Name        Identifier
	// Can not be set together with Body. Optional if Expression is set.
	Type Type
	// Can not be set together with Body. Required if both Type and Body
	// are nil.
	Expression Expression
	// Can not be set together with Type or Expression. Required if both
	// Type and Expression are nil.
	Body *ObjectBody
}

func (p *ClassProperty) isModuleMember() {}
func (p *ClassProperty) isClassMember()  {}

func (p *ClassProperty) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	docs, err := p.Docs.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(docs, newlineWithIndentation(ctx))

	annotations, err := p.Annotations.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(annotations, newlineWithIndentation(ctx))

	modifiers, err := p.Modifiers.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(modifiers, " ")

	name, err := p.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(name)

	if p.Type != nil {
		typ, err := p.Type.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(": ", typ)
	}

	if p.Expression != nil {
		expr, err := p.Expression.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" = ", expr)
	}

	if p.Body != nil {
		body, err := p.Body.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" ", body)
	}

	return b.Bytes(), nil
}

type MethodSignature struct {
	// Optional.
	Modifiers Modifiers
	// Required.
	Name Identifier
	// Optional.
	TypeParameters TypeParameters
	// Optional.
	Parameters Parameters
	// Optional.
	Result Type
}

func (m *MethodSignature) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	modifiers, err := m.Modifiers.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(modifiers, " ")

	name, err := m.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix("function ", name)

	typeParams, err := m.TypeParameters.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(typeParams)

	params, err := joinNodes(ctx, m.Parameters, ", ")
	if err != nil {
		return nil, err
	}
	b.WriteRune('(')
	b.Write(params)
	b.WriteRune(')')

	if m.Result != nil {
		result, err := m.Result.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(": ", result)
	}

	return b.Bytes(), nil
}

// ClassMethod represents a method declaration with an optional implementation
// in a class or module.
type ClassMethod struct {
	// Optional.
	Docs Docs
	// Optional.
	Annotations Annotations
	// Required.
	Signature *MethodSignature
	// Optional.
	Implementation Expression
}

func (m *ClassMethod) isModuleMember() {}
func (m *ClassMethod) isClassMember()  {}

func (m *ClassMethod) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	docs, err := m.Docs.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(docs, newlineWithIndentation(ctx))

	annotations, err := m.Annotations.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(annotations, newlineWithIndentation(ctx))

	signature, err := m.Signature.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(signature)

	if m.Implementation != nil {
		impl, err := m.Implementation.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" = ", impl)
	}

	return b.Bytes(), nil
}
