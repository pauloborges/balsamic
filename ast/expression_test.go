package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node BuiltinExpression
		res  string
		err  error
	}{
		{
			name: "this",
			node: ExpressionThis,
			res:  "this",
		},
		{
			name: "outer",
			node: ExpressionOuter,
			res:  "outer",
		},
		{
			name: "module",
			node: ExpressionModule,
			res:  "module",
		},
		{
			name: "null",
			node: ExpressionNull,
			res:  "null",
		},
		{
			name: "true",
			node: ExpressionTrue,
			res:  "true",
		},
		{
			name: "false",
			node: ExpressionFalse,
			res:  "false",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestIntExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node IntExpression
		res  string
		err  error
	}{
		{
			name: "zero",
			node: 0,
			res:  "0",
		},
		{
			name: "positive",
			node: 42,
			res:  "42",
		},
		{
			name: "negative",
			node: -42,
			res:  "-42",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestFloatExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node FloatExpression
		res  string
		err  error
	}{
		{
			name: "zero",
			node: 0.0,
			res:  "0",
		},
		{
			name: "positive",
			node: 3.1415,
			res:  "3.1415",
		},
		{
			name: "negative",
			node: -3.1415,
			res:  "-3.1415",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestStringExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node StringExpression
		res  string
		err  error
	}{
		{
			name: "empty",
			node: "",
			res:  `""`,
		},
		{
			name: "simple",
			node: "hello",
			res:  `"hello"`,
		},
		{
			name: "with quotes",
			node: `"hello"`,
			res:  `"\"hello\""`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestPrefixUnaryExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node PrefixUnaryExpression
		res  string
		err  error
	}{
		{
			name: "not",
			node: PrefixUnaryExpression{
				Operator: UnaryOperandLogicalNot,
				Operand:  ExpressionTrue,
			},
			res: "!true",
		},
		{
			name: "minus",
			node: PrefixUnaryExpression{
				Operator: UnaryOperandMinus,
				Operand:  IntExpression(42),
			},
			res: "-42",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestPostfixUnaryExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node PostfixUnaryExpression
		res  string
		err  error
	}{
		{
			name: "non-null assertion",
			node: PostfixUnaryExpression{
				Operator: PostfixUnaryOperandNonNullAssertion,
				Operand:  &MemberAccessExpression{Name: "foobar"},
			},
			res: "foobar!!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestBinaryExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node BinaryExpression
		res  string
		err  error
	}{
		{
			name: "addition",
			node: BinaryExpression{
				Operator: BinaryOperatorPlus,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 + 1",
		},
		{
			name: "subtraction",
			node: BinaryExpression{
				Operator: BinaryOperatorMinus,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 - 1",
		},
		{
			name: "multiplication",
			node: BinaryExpression{
				Operator: BinaryOperatorMultiply,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 * 1",
		},
		{
			name: "division",
			node: BinaryExpression{
				Operator: BinaryOperatorDivide,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 / 1",
		},
		{
			name: "integer division",
			node: BinaryExpression{
				Operator: BinaryOperatorIntegerDivide,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 ~/ 1",
		},
		{
			name: "modulo",
			node: BinaryExpression{
				Operator: BinaryOperatorModulo,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 % 1",
		},
		{
			name: "exponentiation",
			node: BinaryExpression{
				Operator: BinaryOperatorExponent,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 ** 1",
		},
		{
			name: "equality",
			node: BinaryExpression{
				Operator: BinaryOperatorEqual,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 == 1",
		},
		{
			name: "not equal",
			node: BinaryExpression{
				Operator: BinaryOperatorNotEqual,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 != 1",
		},
		{
			name: "less than",
			node: BinaryExpression{
				Operator: BinaryOperatorLessThan,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 < 1",
		},
		{
			name: "less than or equal",
			node: BinaryExpression{
				Operator: BinaryOperatorLessThanOrEqual,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 <= 1",
		},
		{
			name: "greater than",
			node: BinaryExpression{
				Operator: BinaryOperatorGreaterThan,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 > 1",
		},
		{
			name: "greater than or equal",
			node: BinaryExpression{
				Operator: BinaryOperatorGreaterThanOrEqual,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 >= 1",
		},
		{
			name: "bitwise and",
			node: BinaryExpression{
				Operator: BinaryOperatorBitwiseAnd,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 & 1",
		},
		{
			name: "bitwise or",
			node: BinaryExpression{
				Operator: BinaryOperatorBitwiseOr,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 | 1",
		},
		{
			name: "logical and",
			node: BinaryExpression{
				Operator: BinaryOperatorLogicalAnd,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 && 1",
		},
		{
			name: "logical or",
			node: BinaryExpression{
				Operator: BinaryOperatorLogicalOr,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 || 1",
		},
		{
			name: "null coalesce",
			node: BinaryExpression{
				Operator: BinaryOperatorNullCoalesce,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 ?? 1",
		},
		{
			name: "pipe",
			node: BinaryExpression{
				Operator: BinaryOperatorPipe,
				Left:     IntExpression(2),
				Right:    IntExpression(1),
			},
			res: "2 |> 1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestTypeExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node TypeExpression
		res  string
		err  error
	}{
		{
			name: "is",
			node: TypeExpression{
				Operator:   TypeOperatorIs,
				Expression: StringExpression("foobar"),
				Type:       &DeclaredType{Name: "String"},
			},
			res: `"foobar" is String`,
		},
		{
			name: "as",
			node: TypeExpression{
				Operator:   TypeOperatorAs,
				Expression: StringExpression("foobar"),
				Type:       &DeclaredType{Name: "String"},
			},
			res: `"foobar" as String`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestMemberAccessExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node MemberAccessExpression
		res  string
		err  error
	}{
		{
			name: "simple",
			node: MemberAccessExpression{Name: "foobar"},
			res:  "foobar",
		},
		{
			name: "call without arguments",
			node: MemberAccessExpression{Name: "foobar", Arguments: NoArguments},
			res:  "foobar()",
		},
		{
			name: "call with arguments",
			node: MemberAccessExpression{
				Name:      "foobar",
				Arguments: Expressions{IntExpression(42), StringExpression("foobar")},
			},
			res: `foobar(42, "foobar")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestQualifiedMemberAccessExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node QualifiedMemberAccessExpression
		res  string
		err  error
	}{
		{
			name: "simple",
			node: QualifiedMemberAccessExpression{
				Receiver: ExpressionThis,
				Name:     "foobar",
			},
			res: "this.foobar",
		},
		{
			name: "nullable",
			node: QualifiedMemberAccessExpression{
				Receiver: ExpressionThis,
				Name:     "foobar",
				Nullable: true,
			},
			res: "this.?foobar",
		},
		{
			name: "call without arguments",
			node: QualifiedMemberAccessExpression{
				Receiver:  ExpressionThis,
				Name:      "foobar",
				Arguments: NoArguments,
			},
			res: "this.foobar()",
		},
		{
			name: "call with arguments",
			node: QualifiedMemberAccessExpression{
				Receiver:  ExpressionThis,
				Name:      "foobar",
				Arguments: Expressions{IntExpression(42), StringExpression("foobar")},
			},
			res: `this.foobar(42, "foobar")`,
		},
		{
			name: "nullable call",
			node: QualifiedMemberAccessExpression{
				Receiver:  ExpressionThis,
				Name:      "foobar",
				Arguments: NoArguments,
				Nullable:  true,
			},
			res: "this.?foobar()",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestSuperAccessExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node SuperAccessExpression
		res  string
		err  error
	}{
		{
			name: "access",
			node: SuperAccessExpression{Name: "foobar"},
			res:  "super.foobar",
		},
		{
			name: "call without arguments",
			node: SuperAccessExpression{Name: "foobar", Arguments: NoArguments},
			res:  "super.foobar()",
		},
		{
			name: "call with arguments",
			node: SuperAccessExpression{
				Name:      "foobar",
				Arguments: Expressions{IntExpression(42), StringExpression("foobar")},
			},
			res: `super.foobar(42, "foobar")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestSubscriptExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node SubscriptExpression
		res  string
		err  error
	}{
		{
			name: "simple",
			node: SubscriptExpression{
				Receiver:  &MemberAccessExpression{Name: "foobar"},
				Subscript: IntExpression(42),
			},
			res: "foobar[42]",
		},
		{
			name: "complex subscript",
			node: SubscriptExpression{
				Receiver:  &MemberAccessExpression{Name: "foobar"},
				Subscript: &MemberAccessExpression{Name: "random", Arguments: NoArguments},
			},
			res: "foobar[random()]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestSuperSubscriptExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node SuperSubscriptExpression
		res  string
		err  error
	}{
		{
			name: "simple",
			node: SuperSubscriptExpression{
				Subscript: IntExpression(42),
			},
			res: "super[42]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestParenthesizedExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ParenthesizedExpression
		res  string
		err  error
	}{
		{
			name: "simple",
			node: ParenthesizedExpression{
				Expression: IntExpression(42),
			},
			res: "(42)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestNewExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node NewExpression
		res  string
		err  error
	}{
		{
			name: "without type",
			node: NewExpression{
				Body: &ObjectBody{},
			},
			res: "new {}",
		},
		{
			name: "with type",
			node: NewExpression{
				Type: &DeclaredType{Name: "Foo"},
				Body: &ObjectBody{},
			},
			res: "new Foo {}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestAmendExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node AmendExpression
		res  string
		err  error
	}{
		{
			name: "amend",
			node: AmendExpression{
				Parent: &ParenthesizedExpression{
					Expression: &MemberAccessExpression{Name: "foobar"},
				},
				Body: &ObjectBody{
					Members: []ObjectMember{
						&ObjectProperty{
							Name:  Identifier("foo"),
							Value: StringExpression("bar"),
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|(foobar) {
				|  foo = "bar"
				|}
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestIfExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node IfExpression
		res  string
		err  error
	}{
		{
			name: "if expression",
			node: IfExpression{
				Condition: &BinaryExpression{
					Operator: BinaryOperatorEqual,
					Left:     &MemberAccessExpression{Name: "answer"},
					Right:    IntExpression(42),
				},
				Then: StringExpression("foo"),
				Else: StringExpression("bar"),
			},
			res: `if (answer == 42) "foo" else "bar"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestImportExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ImportExpression
		res  string
		err  error
	}{
		{
			name: "import",
			node: ImportExpression{
				Path: "./foobar.pkl",
			},
			res: `import("./foobar.pkl")`,
		},
		{
			name: "with glob",
			node: ImportExpression{
				Path: "./*.pkl",
				Glob: true,
			},
			res: `import*("./*.pkl")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestLetExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node LetExpression
		res  string
		err  error
	}{
		{
			name: "let",
			node: LetExpression{
				Name:  &Parameter{Name: Identifier("foo"), Type: &DeclaredType{Name: "String"}},
				Value: StringExpression("bar"),
				Expression: &QualifiedMemberAccessExpression{
					Receiver:  &MemberAccessExpression{Name: "foo"},
					Name:      Identifier("toUpperCase"),
					Arguments: NoArguments,
				},
			},
			res: `let (foo: String = "bar") foo.toUpperCase()`,
		},
		{
			name: "blank",
			node: LetExpression{
				Name:  ParameterBlank,
				Value: StringExpression("ignored"),
				Expression: &ThrowExpression{
					Value: StringExpression("whoops!"),
				},
			},
			res: `let (_ = "ignored") throw("whoops!")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestReadExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ReadExpression
		res  string
		err  error
	}{
		{
			name: "read",
			node: ReadExpression{Value: StringExpression("env:FOOBAR")},
			res:  `read("env:FOOBAR")`,
		},
		{
			name: "nullable",
			node: ReadExpression{
				Value:   StringExpression("env:FOOBAR"),
				Variant: ReadVariantNullable,
			},
			res: `read?("env:FOOBAR")`,
		},
		{
			name: "glob",
			node: ReadExpression{
				Value:   StringExpression("file:*.json"),
				Variant: ReadVariantGlob,
			},
			res: `read*("file:*.json")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestThrowExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node ThrowExpression
		res  string
		err  error
	}{
		{
			name: "throw",
			node: ThrowExpression{Value: StringExpression("error!")},
			res:  `throw("error!")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}

func TestTraceExpressionMarshal(t *testing.T) {
	tests := []struct {
		name string
		node TraceExpression
		res  string
		err  error
	}{
		{
			name: "trace",
			node: TraceExpression{Value: ExpressionThis},
			res:  `trace(this)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.node.Marshal(context.Background())

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.res, string(res))
		})
	}
}
