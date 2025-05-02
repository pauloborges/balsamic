package ast

import (
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Parameter struct {
	// Required.
	Name Identifier
	// Optional.
	Type Type
}

var ParameterBlank = &Parameter{Name: IdentifierBlank}

func (p *Parameter) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

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

	return b.Bytes(), nil
}

type Parameters []*Parameter

func (p Parameters) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, p, ", ")
}

type TypeVariance string

const (
	VarianceIn  TypeVariance = "in"
	VarianceOut TypeVariance = "out"
)

type TypeParameter struct {
	Variance TypeVariance
	Name     Identifier
}

func (t *TypeParameter) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	if t.Variance != "" {
		b.WriteWithSuffix([]byte(t.Variance), " ")
	}

	name, err := t.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.Write(name)

	return b.Bytes(), nil
}

type TypeParameters []*TypeParameter

func (t TypeParameters) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	params, err := joinNodes(ctx, t, ", ")
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("<", params, ">")

	return b.Bytes(), nil
}
