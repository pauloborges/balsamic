package ast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifierMarshal(t *testing.T) {
	tests := []struct {
		name string
		node Identifier
		res  string
		err  error
	}{
		{
			name: "identifier",
			node: "foo",
			res:  "foo",
		},
		{
			name: "empty",
			node: "",
			res:  "",
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

func TestQualifiedIdentifierMarshal(t *testing.T) {
	tests := []struct {
		name string
		node QualifiedIdentifier
		res  string
		err  error
	}{
		{
			name: "qualified identifier",
			node: "foo.bar",
			res:  "foo.bar",
		},
		{
			name: "empty",
			node: "",
			res:  "",
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
