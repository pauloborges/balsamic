package ast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifierMarshal(t *testing.T) {
	tests := []struct {
		name string
		node Modifier
		res  string
		err  error
	}{
		{
			name: "abstract",
			node: ModifierAbstract,
			res:  "abstract",
		},
		{
			name: "const",
			node: ModifierConst,
			res:  "const",
		},
		{
			name: "external",
			node: ModifierExternal,
			res:  "external",
		},
		{
			name: "fixed",
			node: ModifierFixed,
			res:  "fixed",
		},
		{
			name: "hidden",
			node: ModifierHidden,
			res:  "hidden",
		},
		{
			name: "local",
			node: ModifierLocal,
			res:  "local",
		},
		{
			name: "open",
			node: ModifierOpen,
			res:  "open",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}
