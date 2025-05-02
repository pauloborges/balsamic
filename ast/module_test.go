package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestImportClauseTest(t *testing.T) {
	tests := []struct {
		name string
		node ImportClause
		res  string
		err  error
	}{
		{
			name: "import clause",
			node: ImportClause{
				Path: "@foo/Bar.pkl",
			},
			res: `import "@foo/Bar.pkl"`,
		},
		{
			name: "with alias",
			node: ImportClause{
				Path:  "@foo/Bar.pkl",
				Alias: "Baz",
			},
			res: `import "@foo/Bar.pkl" as Baz`,
		},
		{
			name: "with glob",
			node: ImportClause{
				Path:  "@foo/*.pkl",
				Glob:  true,
				Alias: "allFooes",
			},
			res: `import* "@foo/*.pkl" as allFooes`,
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

func TestImportClausesMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ImportClauses
		res  string
		err  error
	}{
		{
			name: "empty",
			node: ImportClauses{},
			res:  "",
		},
		{
			name: "single",
			node: ImportClauses{
				&ImportClause{Path: "@foo/Bar.pkl"},
			},
			res: `import "@foo/Bar.pkl"`,
		},
		{
			name: "several",
			node: ImportClauses{
				&ImportClause{Path: "@foo/Bar.pkl"},
				&ImportClause{Path: "@foo/Baz.pkl"},
			},
			res: stringsutil.StripMargin(`
				|import "@foo/Bar.pkl"
				|import "@foo/Baz.pkl"
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

func TestModuleMembersMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ModuleMembers
		res  string
		err  error
	}{
		{
			name: "empty",
			node: ModuleMembers{},
			res:  "",
		},
		{
			name: "single",
			node: ModuleMembers{
				&ClassProperty{
					Name:       Identifier("foo"),
					Expression: StringExpression("bar"),
				},
			},
			res: `foo = "bar"`,
		},
		{
			name: "several",
			node: ModuleMembers{
				&ClassProperty{
					Name:       Identifier("foo"),
					Expression: StringExpression("bar"),
				},
				&ClassProperty{
					Name:       Identifier("baz"),
					Expression: StringExpression("qux"),
				},
			},
			res: stringsutil.StripMargin(`
				|foo = "bar"
				|
				|baz = "qux"
			`),
		},
		{
			name: "class method is a module member",
			node: ModuleMembers{
				&ClassMethod{
					Signature: &MethodSignature{
						Name:   Identifier("random"),
						Result: &DeclaredType{Name: "Int"},
					},
					Implementation: IntExpression(42),
				},
			},
			res: stringsutil.StripMargin(`
				|function random(): Int = 42
			`),
		},
		{
			name: "class is a module member",
			node: ModuleMembers{
				&Class{
					Name: Identifier("Foo"),
				},
			},
			res: stringsutil.StripMargin(`
				|class Foo
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

func TestModuleMarshal(t *testing.T) {
	tests := []struct {
		name string
		node Module
		res  string
		err  error
	}{
		{
			name: "named",
			node: Module{
				Name: "foo.bar",
			},
			res: stringsutil.StripMargin(`
				|module foo.bar
				|
			`),
		},
		{
			name: "amends",
			node: Module{
				ParentRelationship: ModuleRelationshipAmends,
				ParentName:         "./foo.pkl",
			},
			res: stringsutil.StripMargin(`
				|amends "./foo.pkl"
				|
			`),
		},
		{
			name: "extends",
			node: Module{
				ParentRelationship: ModuleRelationshipExtends,
				ParentName:         "./foo.pkl",
			},
			res: stringsutil.StripMargin(`
				|extends "./foo.pkl"
				|
			`),
		},
		{
			name: "shebang",
			node: Module{
				Name:           "foo.bar",
				ShebangComment: ShebangComment("/usr/bin/pkl eval"),
			},
			res: stringsutil.StripMargin(`
				|#! /usr/bin/pkl eval
				|module foo.bar
				|
			`),
		},
		{
			name: "docs",
			node: Module{
				Name: "foo.bar",
				Docs: Docs("This is a test module"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a test module
				|module foo.bar
				|
			`),
		},
		{
			name: "annotations",
			node: Module{
				Name: "foo.bar",
				Annotations: Annotations{
					&Annotation{Name: "Baz"},
				},
			},
			res: stringsutil.StripMargin(`
				|@Baz
				|module foo.bar
				|
			`),
		},
		{
			name: "modifiers",
			node: Module{
				Name:      "foo.bar",
				Modifiers: Modifiers{ModifierAbstract},
			},
			res: stringsutil.StripMargin(`
				|abstract module foo.bar
				|
			`),
		},
		{
			name: "imports",
			node: Module{
				Name: "foo.bar",
				Imports: ImportClauses{
					&ImportClause{Path: "@foo/Bar.pkl"},
					&ImportClause{Path: "@foo/Baz.pkl"},
				},
			},
			res: stringsutil.StripMargin(`
				|module foo.bar
				|
				|import "@foo/Bar.pkl"
				|import "@foo/Baz.pkl"
				|
			`),
		},
		{
			name: "members",
			node: Module{
				Name: "foo.bar",
				Members: ModuleMembers{
					&ClassProperty{
						Name:       Identifier("foo"),
						Expression: StringExpression("bar"),
					},
					&ClassProperty{
						Name:       Identifier("baz"),
						Expression: StringExpression("qux"),
					},
				},
			},
			res: stringsutil.StripMargin(`
				|module foo.bar
				|
				|foo = "bar"
				|
				|baz = "qux"
				|
			`),
		},
		{
			name: "all together",
			node: Module{
				Name:               "foo.bar",
				ParentRelationship: ModuleRelationshipAmends,
				ParentName:         "./foo.pkl",
				Modifiers:          Modifiers{ModifierAbstract},
				Docs:               Docs("This is a test module"),
				Annotations:        Annotations{&Annotation{Name: "Baz"}},
				ShebangComment:     ShebangComment("/usr/bin/pkl eval"),
				Imports: ImportClauses{
					&ImportClause{Path: "@foo/Bar.pkl"},
					&ImportClause{Path: "@foo/Baz.pkl"},
				},
				Members: ModuleMembers{
					&ClassProperty{
						Name:       Identifier("foo"),
						Expression: StringExpression("bar"),
					},
					&ClassProperty{
						Name:       Identifier("baz"),
						Expression: StringExpression("qux"),
					},
				},
			},
			res: stringsutil.StripMargin(`
				|#! /usr/bin/pkl eval
				|/// This is a test module
				|@Baz
				|abstract module foo.bar amends "./foo.pkl"
				|
				|import "@foo/Bar.pkl"
				|import "@foo/Baz.pkl"
				|
				|foo = "bar"
				|
				|baz = "qux"
				|
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
