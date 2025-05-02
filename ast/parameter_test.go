package ast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterMarshal(t *testing.T) {
	tests := []struct {
		name string
		node Parameter
		res  string
		err  error
	}{
		{
			name: "without type",
			node: Parameter{Name: "name"},
			res:  "name",
		},
		{
			name: "with type",
			node: Parameter{
				Name: "name",
				Type: &DeclaredType{Name: "String"},
			},
			res: "name: String",
		},
		{
			name: "blank",
			node: *ParameterBlank,
			res:  "_",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.node.Marshal(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}

func TestTypeParameterMarshal(t *testing.T) {
	tests := []struct {
		name string
		node TypeParameter
		res  string
		err  error
	}{
		{
			name: "type parameter",
			node: TypeParameter{Name: "T"},
			res:  "T",
		},
		{
			name: "in",
			node: TypeParameter{
				Name:     "T",
				Variance: VarianceIn,
			},
			res: "in T",
		},
		{
			name: "out",
			node: TypeParameter{
				Name:     "T",
				Variance: VarianceOut,
			},
			res: "out T",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.node.Marshal(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}
