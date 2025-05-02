package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestObjectBodyMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectBody
		res  string
		err  error
	}{
		{
			name: "empty",
			node: ObjectBody{},
			res:  "{}",
		},
		{
			name: "with parameters",
			node: ObjectBody{
				Parameters: Parameters{
					&Parameter{Name: "foo", Type: &DeclaredType{Name: "String"}},
					&Parameter{Name: "bar", Type: &DeclaredType{Name: "Int"}},
				},
			},
			res: "{ foo: String, bar: Int ->}",
		},
		{
			name: "with members",
			node: ObjectBody{
				Members: ObjectMembers{
					&ObjectProperty{
						Name:  "foo",
						Value: StringExpression("bar"),
					},
					&ObjectProperty{
						Name:  "baz",
						Value: IntExpression(42),
					},
				},
			},
			res: stringsutil.StripMargin(`
				|{
				|  foo = "bar"
				|  baz = 42
				|}
			`),
		},
		{
			name: "with parameters and members",
			node: ObjectBody{
				Parameters: Parameters{
					&Parameter{Name: "param", Type: &DeclaredType{Name: "String"}},
				},
				Members: ObjectMembers{
					&ObjectProperty{
						Name:  "foo",
						Value: StringExpression("bar"),
					},
				},
			},
			res: stringsutil.StripMargin(`
				|{ param: String ->
				|  foo = "bar"
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

func TestObjectMembersMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectMembers
		res  string
		err  error
	}{
		{
			name: "empty",
			node: ObjectMembers{},
			res:  "",
		},
		{
			name: "single property",
			node: ObjectMembers{
				&ObjectProperty{
					Name:  "foo",
					Value: StringExpression("bar"),
				},
			},
			res: stringsutil.StripMargin(`
				|
				|  foo = "bar"
			`),
		},
		{
			name: "multiple properties",
			node: ObjectMembers{
				&ObjectProperty{
					Name:  "foo",
					Value: StringExpression("bar"),
				},
				&ObjectProperty{
					Name:  "baz",
					Value: IntExpression(42),
				},
			},
			res: stringsutil.StripMargin(`
				|
				|  foo = "bar"
				|  baz = 42
			`),
		},
		{
			name: "mixed members",
			node: ObjectMembers{
				&ObjectProperty{
					Name:  "foo",
					Value: StringExpression("bar"),
				},
				&ObjectElement{
					Value: IntExpression(42),
				},
			},
			res: stringsutil.StripMargin(`
				|
				|  foo = "bar"
				|  42
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

func TestObjectPropertyMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectProperty
		res  string
		err  error
	}{
		{
			name: "with value only",
			node: ObjectProperty{
				Name:  "foo",
				Value: StringExpression("bar"),
			},
			res: `foo = "bar"`,
		},
		{
			name: "with type only",
			node: ObjectProperty{
				Name: "foo",
				Type: &DeclaredType{Name: "String"},
			},
			res: `foo: String`,
		},
		{
			name: "with type and value",
			node: ObjectProperty{
				Name:  "foo",
				Type:  &DeclaredType{Name: "String"},
				Value: StringExpression("bar"),
			},
			res: `foo: String = "bar"`,
		},
		{
			name: "with body",
			node: ObjectProperty{
				Name: "foo",
				Body: []*ObjectBody{
					{
						Members: ObjectMembers{
							&ObjectProperty{
								Name:  "bar",
								Value: StringExpression("baz"),
							},
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
			name: "with modifiers",
			node: ObjectProperty{
				Modifiers: Modifiers{ModifierConst},
				Name:      "foo",
				Value:     StringExpression("bar"),
			},
			res: `const foo = "bar"`,
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

func TestObjectMethodMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectMethod
		res  string
		err  error
	}{
		{
			name: "simple",
			node: ObjectMethod{
				Signature: &MethodSignature{
					Name: "foo",
				},
				Value: StringExpression("bar"),
			},
			res: `function foo() = "bar"`,
		},
		{
			name: "with parameters",
			node: ObjectMethod{
				Signature: &MethodSignature{
					Name: "foo",
					Parameters: Parameters{
						&Parameter{Name: "param", Type: &DeclaredType{Name: "String"}},
					},
				},
				Value: &MemberAccessExpression{Name: "param"},
			},
			res: `function foo(param: String) = param`,
		},
		{
			name: "with return type",
			node: ObjectMethod{
				Signature: &MethodSignature{
					Name:   "foo",
					Result: &DeclaredType{Name: "String"},
				},
				Value: StringExpression("bar"),
			},
			res: `function foo(): String = "bar"`,
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

func TestObjectEntryMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectEntry
		res  string
		err  error
	}{
		{
			name: "with value",
			node: ObjectEntry{
				Key:   StringExpression("foo"),
				Value: StringExpression("bar"),
			},
			res: `["foo"] = "bar"`,
		},
		{
			name: "with expression key",
			node: ObjectEntry{
				Key:   &MemberAccessExpression{Name: "foo", Arguments: NoArguments},
				Value: StringExpression("bar"),
			},
			res: `[foo()] = "bar"`,
		},
		{
			name: "with body",
			node: ObjectEntry{
				Key: StringExpression("foo"),
				Body: []*ObjectBody{
					{
						Members: ObjectMembers{
							&ObjectProperty{
								Name:  "bar",
								Value: StringExpression("baz"),
							},
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|["foo"] {
				|  bar = "baz"
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

func TestObjectElementMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectElement
		res  string
		err  error
	}{
		{
			name: "string",
			node: ObjectElement{
				Value: StringExpression("foo"),
			},
			res: `"foo"`,
		},
		{
			name: "integer",
			node: ObjectElement{
				Value: IntExpression(42),
			},
			res: `42`,
		},
		{
			name: "expression",
			node: ObjectElement{
				Value: &BinaryExpression{
					Operator: BinaryOperatorPlus,
					Left:     IntExpression(1),
					Right:    IntExpression(2),
				},
			},
			res: `1 + 2`,
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

func TestObjectSpreadMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ObjectSpread
		res  string
		err  error
	}{
		{
			name: "regular spread",
			node: ObjectSpread{
				Value: &MemberAccessExpression{Name: "foo"},
			},
			res: `...foo`,
		},
		{
			name: "nullable spread",
			node: ObjectSpread{
				Value:    &MemberAccessExpression{Name: "foo"},
				Nullable: true,
			},
			res: `...?foo`,
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

func TestMemberPredicateMarshal(t *testing.T) {
	tests := []struct {
		name string
		node MemberPredicate
		res  string
		err  error
	}{
		{
			name: "with value",
			node: MemberPredicate{
				Condition: &BinaryExpression{
					Operator: BinaryOperatorEqual,
					Left:     &MemberAccessExpression{Name: "foo"},
					Right:    StringExpression("bar"),
				},
				Value: StringExpression("baz"),
			},
			res: `[[foo == "bar"]] = "baz"`,
		},
		{
			name: "with body",
			node: MemberPredicate{
				Condition: &BinaryExpression{
					Operator: BinaryOperatorEqual,
					Left:     &MemberAccessExpression{Name: "foo"},
					Right:    StringExpression("bar"),
				},
				Body: []*ObjectBody{
					{
						Members: ObjectMembers{
							&ObjectProperty{
								Name:  "qux",
								Value: StringExpression("quux"),
							},
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|[[foo == "bar"]] {
				|  qux = "quux"
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

func TestForGeneratorMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ForGenerator
		res  string
		err  error
	}{
		{
			name: "simple",
			node: ForGenerator{
				Value:      &Parameter{Name: "item"},
				Collection: &MemberAccessExpression{Name: "items"},
				Body: &ObjectBody{
					Members: ObjectMembers{
						&ObjectProperty{
							Name:  "value",
							Value: &MemberAccessExpression{Name: "item"},
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|for (item in items) {
				|  value = item
				|}
			`),
		},
		{
			name: "with key",
			node: ForGenerator{
				Key:        &Parameter{Name: "index"},
				Value:      &Parameter{Name: "item"},
				Collection: &MemberAccessExpression{Name: "items"},
				Body: &ObjectBody{
					Members: ObjectMembers{
						&ObjectProperty{
							Name:  "index",
							Value: &MemberAccessExpression{Name: "index"},
						},
						&ObjectProperty{
							Name:  "value",
							Value: &MemberAccessExpression{Name: "item"},
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|for (index, item in items) {
				|  index = index
				|  value = item
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

func TestWhenGeneratorMarshal(t *testing.T) {
	tests := []struct {
		name string
		node WhenGenerator
		res  string
		err  error
	}{
		{
			name: "without else",
			node: WhenGenerator{
				Condition: &BinaryExpression{
					Operator: BinaryOperatorEqual,
					Left:     &MemberAccessExpression{Name: "foo"},
					Right:    StringExpression("bar"),
				},
				Then: &ObjectBody{
					Members: ObjectMembers{
						&ObjectProperty{
							Name:  "value",
							Value: StringExpression("baz"),
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|when (foo == "bar") {
				|  value = "baz"
				|}
			`),
		},
		{
			name: "with else",
			node: WhenGenerator{
				Condition: &BinaryExpression{
					Operator: BinaryOperatorEqual,
					Left:     &MemberAccessExpression{Name: "foo"},
					Right:    StringExpression("bar"),
				},
				Then: &ObjectBody{
					Members: ObjectMembers{
						&ObjectProperty{
							Name:  "value",
							Value: StringExpression("baz"),
						},
					},
				},
				Else: &ObjectBody{
					Members: ObjectMembers{
						&ObjectProperty{
							Name:  "value",
							Value: StringExpression("qux"),
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|when (foo == "bar") {
				|  value = "baz"
				|} else {
				|  value = "qux"
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
