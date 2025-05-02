package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestTypeAliasMarshal(t *testing.T) {
	tests := []struct {
		name string
		node TypeAlias
		res  string
		err  error
	}{
		{
			name: "type alias",
			node: TypeAlias{
				Name: Identifier("HelloWorld"),
				Type: StringLiteralType("helloworld"),
			},
			res: `typealias HelloWorld = "helloworld"`,
		},
		{
			name: "with docs",
			node: TypeAlias{
				Docs: Docs("This is a Hello World.\n\nMore docs here."),
				Name: Identifier("HelloWorld"),
				Type: StringLiteralType("helloworld"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a Hello World.
				|///
				|/// More docs here.
				|typealias HelloWorld = "helloworld"
			`),
		},
		{
			name: "with annotations",
			node: TypeAlias{
				Annotations: Annotations{
					&Annotation{Name: "Foo"},
					&Annotation{Name: "Bar"},
				},
				Name: Identifier("HelloWorld"),
				Type: StringLiteralType("helloworld"),
			},
			res: stringsutil.StripMargin(`
				|@Foo
				|@Bar
				|typealias HelloWorld = "helloworld"
			`),
		},
		{
			name: "with modifiers",
			node: TypeAlias{
				Modifiers: Modifiers{ModifierLocal, ModifierConst},
				Name:      Identifier("HelloWorld"),
				Type:      StringLiteralType("helloworld"),
			},
			res: stringsutil.StripMargin(`local const typealias HelloWorld = "helloworld"`),
		},
		{
			name: "with parameters",
			node: TypeAlias{
				Name: Identifier("HelloWorld"),
				Parameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")}},
				Type: StringLiteralType("helloworld"),
			},
			res: stringsutil.StripMargin(`typealias HelloWorld<T> = "helloworld"`),
		},
		{
			name: "with all",
			node: TypeAlias{
				Docs: Docs("This is a Hello World.\n\nMore docs here."),
				Annotations: Annotations{
					&Annotation{Name: "Foo"},
					&Annotation{Name: "Bar"},
				},
				Modifiers: Modifiers{ModifierLocal, ModifierConst},
				Name:      Identifier("HelloWorld"),
				Parameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")}},
				Type: StringLiteralType("helloworld"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a Hello World.
				|///
				|/// More docs here.
				|@Foo
				|@Bar
				|local const typealias HelloWorld<T> = "helloworld"
			`),
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
