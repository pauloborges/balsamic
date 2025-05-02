package ast

import "context"

type Modifier string

const (
	ModifierAbstract Modifier = "abstract"
	ModifierConst    Modifier = "const"
	ModifierExternal Modifier = "external"
	ModifierFixed    Modifier = "fixed"
	ModifierHidden   Modifier = "hidden"
	ModifierLocal    Modifier = "local"
	ModifierOpen     Modifier = "open"
)

func (m Modifier) Marshal(_ context.Context) ([]byte, error) {
	return []byte(m), nil
}

type Modifiers []Modifier

func (m Modifiers) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, m, " ")
}
