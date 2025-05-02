package ast

type PrefixUnaryOperand string

const (
	UnaryOperandMinus      PrefixUnaryOperand = "-"
	UnaryOperandLogicalNot PrefixUnaryOperand = "!"
)

type PostfixUnaryOperand string

const (
	PostfixUnaryOperandNonNullAssertion PostfixUnaryOperand = "!!"
)

type BinaryOperator string

const (
	BinaryOperatorPlus               BinaryOperator = "+"
	BinaryOperatorMinus              BinaryOperator = "-"
	BinaryOperatorMultiply           BinaryOperator = "*"
	BinaryOperatorDivide             BinaryOperator = "/"
	BinaryOperatorIntegerDivide      BinaryOperator = "~/"
	BinaryOperatorModulo             BinaryOperator = "%"
	BinaryOperatorExponent           BinaryOperator = "**"
	BinaryOperatorEqual              BinaryOperator = "=="
	BinaryOperatorNotEqual           BinaryOperator = "!="
	BinaryOperatorLessThan           BinaryOperator = "<"
	BinaryOperatorLessThanOrEqual    BinaryOperator = "<="
	BinaryOperatorGreaterThan        BinaryOperator = ">"
	BinaryOperatorGreaterThanOrEqual BinaryOperator = ">="
	BinaryOperatorBitwiseAnd         BinaryOperator = "&"
	BinaryOperatorBitwiseOr          BinaryOperator = "|"
	BinaryOperatorLogicalAnd         BinaryOperator = "&&"
	BinaryOperatorLogicalOr          BinaryOperator = "||"
	BinaryOperatorNullCoalesce       BinaryOperator = "??"
	BinaryOperatorPipe               BinaryOperator = "|>"
)

type TypeOperator string

const (
	TypeOperatorIs TypeOperator = "is"
	TypeOperatorAs TypeOperator = "as"
)
