package ast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node BuiltinType
		res  string
		err  error
	}{
		{
			name: "unknown",
			node: TypeUnknown,
			res:  "unknown",
		},
		{
			name: "nothing",
			node: TypeNothing,
			res:  "nothing",
		},
		{
			name: "module",
			node: TypeModule,
			res:  "module",
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

func TestStringLiteralTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node StringLiteralType
		res  string
		err  error
	}{
		{
			name: "string literal",
			node: StringLiteralType("hello"),
			res:  `"hello"`,
		},
		{
			name: "empty string literal",
			node: StringLiteralType(""),
			res:  `""`,
		},
		{
			name: "with quotes",
			node: StringLiteralType(`"hello"`),
			res:  `"\"hello\""`,
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

func TestDeclaredTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node DeclaredType
		res  string
		err  error
	}{
		{
			name: "declared type",
			node: DeclaredType{Name: QualifiedIdentifier("foo.Bar")},
			res:  "foo.Bar",
		},
		{
			name: "type arguments",
			node: DeclaredType{
				Name: QualifiedIdentifier("foo.Bar"),
				TypeParameters: []Type{
					&DeclaredType{Name: "String"},
					TypeUnknown,
				},
			},
			res: `foo.Bar<String, unknown>`,
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

func TestParenthesizedTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ParenthesizedType
		res  string
		err  error
	}{
		{
			name: "parenthesized type",
			node: ParenthesizedType{
				Type: &DeclaredType{Name: "String"},
			},
			res: "(String)",
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

func TestNullableTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node NullableType
		res  string
		err  error
	}{
		{
			name: "nullable type",
			node: NullableType{
				Type: &DeclaredType{Name: "String"},
			},
			res: "String?",
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

func TestConstrainedTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ConstrainedType
		res  string
		err  error
	}{
		{
			name: "constrained type",
			node: ConstrainedType{
				Type: &DeclaredType{"String", nil},
				Constraints: Expressions{
					&BinaryExpression{
						Operator: BinaryOperatorLessThanOrEqual,
						Left:     &MemberAccessExpression{Name: "length"},
						Right:    IntExpression(42),
					},
				},
			},
			res: `String(length <= 42)`,
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

func TestUnionTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node *UnionType
		res  string
		err  error
	}{
		{
			name: "union type",
			node: &UnionType{
				Members: []Type{
					&DeclaredType{Name: "Int"},
					StringLiteralType("foobar"),
				},
			},
			res: `Int | "foobar"`,
		},

		{
			name: "default",
			node: &UnionType{
				Members: []Type{
					&DeclaredType{Name: "Int"},
					&DeclaredType{Name: "String"},
				},
				Default: StringLiteralType("foobar"),
			},
			res: `Int | String | *"foobar"`,
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

func TestFunctionLiteralTypeMarshal(t *testing.T) {
	tests := []struct {
		name string
		node FunctionLiteralType
		res  string
		err  error
	}{
		{
			name: "function literal",
			node: FunctionLiteralType{
				Parameters: []Type{
					&DeclaredType{Name: "String"},
					&DeclaredType{Name: "Int"},
				},
				Result: &DeclaredType{Name: "Boolean"},
			},
			res: `(String, Int) -> Boolean`,
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
