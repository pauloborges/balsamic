package ast

import (
	"context"
)

const IdentifierBlank Identifier = "_"

type Identifier string

func (i Identifier) Marshal(_ context.Context) ([]byte, error) {
	return []byte(i), nil
}

type QualifiedIdentifier string

func (i QualifiedIdentifier) Marshal(_ context.Context) ([]byte, error) {
	return []byte(i), nil
}
