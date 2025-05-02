package ast

import (
	"bytes"
	"context"
	"strings"
)

// Node represents a node in Pkl's abstract syntax tree.
type Node interface {
	// Marshal pretty-prints the node and returns the result as a byte slice.
	Marshal(ctx context.Context) ([]byte, error)
}

type ctxKey string

var indentLevelKey ctxKey = "indentLevel"

func ctxWithIndentLevel(ctx context.Context, indentLevel uint) context.Context {
	return context.WithValue(ctx, indentLevelKey, indentLevel)
}

func raiseIndentLevel(ctx context.Context) context.Context {
	indentLevel := getIndentLevel(ctx)
	return context.WithValue(ctx, indentLevelKey, indentLevel+1)
}

func getIndentLevel(ctx context.Context) uint {
	if indent, ok := ctx.Value(indentLevelKey).(uint); ok {
		return indent
	}
	return 0
}

func indentation(ctx context.Context) string {
	indentLevel := getIndentLevel(ctx)
	if indentLevel == 0 {
		return ""
	}
	return strings.Repeat("  ", int(indentLevel))
}

func newlineWithIndentation(ctx context.Context) string {
	return "\n" + indentation(ctx)
}

func joinNodes[N Node](ctx context.Context, nodes []N, separator string) ([]byte, error) {
	var b bytes.Buffer

	if len(nodes) == 0 {
		return nil, nil
	}

	for i, node := range nodes {
		if i > 0 {
			b.WriteString(separator)
		}

		node, err := node.Marshal(ctx)
		if err != nil {
			return nil, err
		}

		b.Write(node)
	}

	return b.Bytes(), nil
}
