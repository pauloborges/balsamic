package ast

import (
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Annotation struct {
	// Name of the class. Required.
	Name QualifiedIdentifier

	// Body of the annotation. Optional.
	Body *ObjectBody
}

func (a *Annotation) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	name, err := a.Name.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefix("@", name)

	if a.Body != nil {
		body, err := a.Body.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix(" ", body)
	}

	return b.Bytes(), nil
}

type Annotations []*Annotation

func (a Annotations) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, a, newlineWithIndentation(ctx))
}
