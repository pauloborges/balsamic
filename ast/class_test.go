package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestClassPropertyMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ClassProperty
		res  string
		err  error
	}{
		{
			name: "with type only",
			node: ClassProperty{
				Name: Identifier("foo"),
				Type: &DeclaredType{Name: "String"},
			},
			res: `foo: String`,
		},
		{
			name: "with expression only",
			node: ClassProperty{
				Name:       Identifier("foo"),
				Expression: StringExpression("bar"),
			},
			res: `foo = "bar"`,
		},
		{
			name: "with both type and expression",
			node: ClassProperty{
				Name:       Identifier("foo"),
				Type:       &DeclaredType{Name: "String"},
				Expression: StringExpression("bar"),
			},
			res: `foo: String = "bar"`,
		},
		{
			name: "with body",
			node: ClassProperty{
				Name: Identifier("foo"),
				Body: &ObjectBody{
					Members: []ObjectMember{
						&ObjectProperty{
							Name:  Identifier("bar"),
							Value: StringExpression("baz"),
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|foo {
				|  bar = "baz"
				|}
			`),
		},
		{
			name: "with docs, annotations and modifiers",
			node: ClassProperty{
				Docs: Docs("This is a foo property."),
				Annotations: Annotations{
					&Annotation{Name: "Foo"},
				},
				Modifiers:  Modifiers{ModifierConst},
				Name:       Identifier("foo"),
				Type:       &DeclaredType{Name: "String"},
				Expression: StringExpression("bar"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a foo property.
				|@Foo
				|const foo: String = "bar"
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

func TestMethodSignatureMarshal(t *testing.T) {
	tests := []struct {
		name string
		node MethodSignature
		res  string
		err  error
	}{
		{
			name: "simple",
			node: MethodSignature{
				Name: Identifier("foo"),
			},
			res: `function foo()`,
		},
		{
			name: "with result type",
			node: MethodSignature{
				Name:   Identifier("foo"),
				Result: &DeclaredType{Name: "String"},
			},
			res: `function foo(): String`,
		},
		{
			name: "with params",
			node: MethodSignature{
				Name: Identifier("foo"),
				Parameters: Parameters{
					&Parameter{Name: Identifier("bar"), Type: &DeclaredType{Name: "String"}},
					&Parameter{Name: Identifier("baz"), Type: &DeclaredType{Name: "Int"}},
				},
			},
			res: `function foo(bar: String, baz: Int)`,
		},
		{
			name: "with type parameters",
			node: MethodSignature{
				Name: Identifier("foo"),
				TypeParameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")},
				},
				Parameters: Parameters{
					&Parameter{Name: Identifier("bar"), Type: &DeclaredType{Name: "T"}},
				},
			},
			res: `function foo<T>(bar: T)`,
		},
		{
			name: "with modifiers",
			node: MethodSignature{
				Modifiers: Modifiers{ModifierAbstract},
				Name:      Identifier("foo"),
			},
			res: `abstract function foo()`,
		},
		{
			name: "with all",
			node: MethodSignature{
				Modifiers: Modifiers{ModifierAbstract},
				Name:      Identifier("foo"),
				TypeParameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")},
					&TypeParameter{Name: Identifier("U")},
				},
				Parameters: Parameters{
					&Parameter{Name: Identifier("bar"), Type: &DeclaredType{Name: "T"}},
				},
				Result: &DeclaredType{Name: "U"},
			},
			res: `abstract function foo<T, U>(bar: T): U`,
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

func TestClassMethodMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ClassMethod
		res  string
		err  error
	}{
		{
			name: "simple",
			node: ClassMethod{
				Signature: &MethodSignature{
					Name: Identifier("foo"),
				},
			},
			res: `function foo()`,
		},
		{
			name: "with implementation",
			node: ClassMethod{
				Signature: &MethodSignature{
					Name: Identifier("foo"),
				},
				Implementation: StringExpression("bar"),
			},
			res: `function foo() = "bar"`,
		},
		{
			name: "with docs",
			node: ClassMethod{
				Docs: Docs("This is a foo method."),
				Signature: &MethodSignature{
					Name: Identifier("foo"),
				},
			},
			res: stringsutil.StripMargin(`
				|/// This is a foo method.
				|function foo()
			`),
		},
		{
			name: "with annotations",
			node: ClassMethod{
				Annotations: Annotations{
					&Annotation{Name: "Foo"},
				},
				Signature: &MethodSignature{
					Name: Identifier("foo"),
				},
			},
			res: stringsutil.StripMargin(`
				|@Foo
				|function foo()
			`),
		},
		{
			name: "with all",
			node: ClassMethod{
				Docs: Docs("This is a foo method."),
				Annotations: Annotations{
					&Annotation{Name: "Foo"},
				},
				Signature: &MethodSignature{
					Name: Identifier("foo"),
				},
				Implementation: StringExpression("bar"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a foo method.
				|@Foo
				|function foo() = "bar"
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

func TestClassMarshal(t *testing.T) {
	tests := []struct {
		name string
		node Class
		res  string
		err  error
	}{
		{
			name: "simple",
			node: Class{
				Name: Identifier("Foo"),
			},
			res: `class Foo`,
		},
		{
			name: "with docs",
			node: Class{
				Docs: Docs("This is a Foo class."),
				Name: Identifier("Foo"),
			},
			res: stringsutil.StripMargin(`
				|/// This is a Foo class.
				|class Foo
			`),
		},
		{
			name: "with annotations",
			node: Class{
				Annotations: Annotations{
					&Annotation{Name: "Bar"},
				},
				Name: Identifier("Foo"),
			},
			res: stringsutil.StripMargin(`
				|@Bar
				|class Foo
			`),
		},
		{
			name: "with modifiers",
			node: Class{
				Modifiers: Modifiers{ModifierAbstract},
				Name:      Identifier("Foo"),
			},
			res: `abstract class Foo`,
		},
		{
			name: "with type parameters",
			node: Class{
				Name: Identifier("Foo"),
				TypeParameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")},
					&TypeParameter{Name: Identifier("U")},
				},
			},
			res: `class Foo<T, U>`,
		},
		{
			name: "with parent",
			node: Class{
				Name:       Identifier("Foo"),
				ParentName: QualifiedIdentifier("Bar"),
			},
			res: `class Foo extends Bar`,
		},
		{
			name: "with parent and type parameters",
			node: Class{
				Name:       Identifier("Foo"),
				ParentName: QualifiedIdentifier("Bar"),
				ParentTypeParameters: TypeParameters{
					&TypeParameter{Name: Identifier("T")},
				},
			},
			res: `class Foo extends Bar<T>`,
		},
		{
			name: "with members",
			node: Class{
				Name: Identifier("Foo"),
				Members: []ClassMember{
					&ClassProperty{
						Name: Identifier("bar"),
						Type: &DeclaredType{Name: "String"},
					},
					&ClassMethod{
						Signature: &MethodSignature{
							Name:   Identifier("baz"),
							Result: &DeclaredType{Name: "Int"},
						},
						Implementation: IntExpression(42),
					},
				},
			},
			res: stringsutil.StripMargin(`
				|class Foo {
				|  bar: String
				|
				|  function baz(): Int = 42
				|}
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
