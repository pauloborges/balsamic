package ast

import (
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type ObjectBody struct {
	Parameters Parameters
	Members    ObjectMembers
}

func (o *ObjectBody) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	b.WriteRune('{')

	params, err := joinNodes(ctx, o.Parameters, ", ")
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix(" ", params, " ->")

	members, err := o.Members.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(members, newlineWithIndentation(ctx))

	b.WriteRune('}')

	return b.Bytes(), nil
}

type ObjectMember interface {
	Node
	isObjectMember()
}

type ObjectMembers []ObjectMember

func (m ObjectMembers) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	raisedCtx := raiseIndentLevel(ctx)
	raisedNl := newlineWithIndentation(raisedCtx)
	members, err := joinNodes(raisedCtx, m, raisedNl)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(raisedNl, members)

	return b.Bytes(), nil
}

type ObjectProperty struct {
	// Optional.
	Modifiers Modifiers
	// Required.
	Name Identifier
	// Can not be set together with Body. Optional if Value is set.
	Type Type
	// Can not be set together with Body. Required if Body is not set.
	Value Expression
	// Can not be set together with Value. Required if Value is not set.
	Body []*ObjectBody
}

func (m *ObjectProperty) isObjectMember() {}

func (p *ObjectProperty) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

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

	if p.Value != nil {
		val, err := p.Value.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" = ", val)
	}

	body, err := joinNodes(ctx, p.Body, " ")
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" ", body)

	return b.Bytes(), nil
}

type ObjectMethod struct {
	// Required.
	Signature *MethodSignature
	// Required.
	Value Expression
}

func (m *ObjectMethod) isObjectMember() {}

func (m *ObjectMethod) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	signature, err := m.Signature.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(signature)

	val, err := m.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" = ", val)

	return b.Bytes(), nil
}

type ObjectEntry struct {
	// Required.
	Key Expression
	// Required if Body is unset. Can not be set together with Body.
	Value Expression
	// Required if Value is unset. Can not be set together with Value.
	Body []*ObjectBody
}

func (m *ObjectEntry) isObjectMember() {}

func (m *ObjectEntry) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	key, err := m.Key.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("[", key, "]")

	if m.Value != nil {
		val, err := m.Value.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" = ", val)
	}

	if m.Body != nil {
		body, err := joinNodes(ctx, m.Body, " ")
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" ", body)
	}

	return b.Bytes(), nil
}

type ObjectElement struct {
	// Required.
	Value Expression
}

func (m *ObjectElement) isObjectMember() {}

func (m *ObjectElement) Marshal(ctx context.Context) ([]byte, error) {
	return m.Value.Marshal(ctx)
}

type ObjectSpread struct {
	// Required.
	Value Expression
	// Optional.
	Nullable bool
}

func (m *ObjectSpread) isObjectMember() {}

func (m *ObjectSpread) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	if m.Nullable {
		b.WriteString("...?")
	} else {
		b.WriteString("...")
	}

	val, err := m.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(val)

	return b.Bytes(), nil
}

type MemberPredicate struct {
	// Required.
	Condition Expression
	// Required if Body is unset. Can not be set together with Body.
	Value Expression
	// Required if Value is unset. Can not be set together with Value.
	Body []*ObjectBody
}

func (m *MemberPredicate) isObjectMember() {}

func (m *MemberPredicate) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer
	cond, err := m.Condition.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("[[", cond, "]]")

	if m.Value != nil {
		val, err := m.Value.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" = ", val)
	}

	if m.Body != nil {
		body, err := joinNodes(ctx, m.Body, " ")
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" ", body)
	}

	return b.Bytes(), nil
}

type ForGenerator struct {
	// Optional.
	Key *Parameter
	// Required.
	Value *Parameter
	// Required.
	Collection Expression
	// Required.
	Body *ObjectBody
}

func (m *ForGenerator) isObjectMember() {}

func (m *ForGenerator) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	b.WriteString("for (")

	if m.Key != nil {
		key, err := m.Key.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithSuffix(key, ", ")
	}

	val, err := m.Value.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(val)

	col, err := m.Collection.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix(" in ", col, ") ")

	body, err := m.Body.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(body)

	return b.Bytes(), nil
}

type WhenGenerator struct {
	// Required.
	Condition Expression
	// Required.
	Then *ObjectBody
	// Optional.
	Else *ObjectBody
}

func (m *WhenGenerator) isObjectMember() {}

func (m *WhenGenerator) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	cond, err := m.Condition.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("when (", cond, ") ")

	thenExpr, err := m.Then.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(thenExpr)

	if m.Else != nil {
		elseExpr, err := m.Else.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" else ", elseExpr)
	}

	return b.Bytes(), nil
}
