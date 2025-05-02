package ast

import (
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type TypeAlias struct {
	// Optional.
	Docs Docs
	// Optional.
	Annotations Annotations
	// Optional.
	Modifiers Modifiers
	// Required.
	Name Identifier
	// Optional.
	Parameters TypeParameters
	// Required.
	Type Type
}

func (t *TypeAlias) isModuleMember() {}

func (t *TypeAlias) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	docs, err := t.Docs.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(docs, newlineWithIndentation(ctx))

	annotations, err := t.Annotations.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(annotations, newlineWithIndentation(ctx))

	modifiers, err := t.Modifiers.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(modifiers, " ")

	name, err := t.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix("typealias ", name)

	params, err := t.Parameters.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(params)

	typ, err := t.Type.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix(" = ", typ)

	return b.Bytes(), nil
}
