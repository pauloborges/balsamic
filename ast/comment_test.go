package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestLineCommentMarshal(t *testing.T) {
	tests := []struct {
		name        string
		indentLevel uint
		node        LineComment
		res         string
		err         error
	}{
		{
			name: "empty",
			node: "",
			res:  "",
		},
		{
			name: "single line",
			node: "This is a comment",
			res:  "// This is a comment",
		},
		{
			name:        "indented",
			indentLevel: 1,
			node:        "This is a comment",
			res:         "  // This is a comment",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			ctx := ctxWithIndentLevel(context.Background(), test.indentLevel)

			// When
			result, err := test.node.Marshal(ctx)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}

func TestBlockCommentMarshal(t *testing.T) {
	tests := []struct {
		name        string
		indentLevel uint
		node        BlockComment
		res         string
		err         error
	}{
		{
			name: "empty",
			node: "",
			res:  "",
		},
		{
			name: "single line",
			node: "This is a comment",
			res: stringsutil.StripMargin(`
				|/*
				|  This is a comment
				|*/
			`),
		},
		{
			name: "multi line",
			node: "This is a comment\nThis is another comment",
			res: stringsutil.StripMargin(`
				|/*
				|  This is a comment
				|  This is another comment
				|*/
			`),
		},
		{
			name: "multi line with empty line",
			node: "This is a comment\n\nThis is another comment",
			res: stringsutil.StripMargin(`
				|/*
				|  This is a comment
				|
				|  This is another comment
				|*/
			`),
		},
		{
			name:        "indented",
			indentLevel: 1,
			node:        "This is a comment",
			res: stringsutil.StripMargin(`
				|/*
				|    This is a comment
				|  */
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			ctx := ctxWithIndentLevel(context.Background(), test.indentLevel)

			// When
			result, err := test.node.Marshal(ctx)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}

func TestDocstMarshal(t *testing.T) {
	tests := []struct {
		name        string
		indentLevel uint
		node        Docs
		res         string
		err         error
	}{
		{
			name: "empty",
			node: "",
			res:  "",
		},
		{
			name: "single line",
			node: "This is a comment",
			res:  "/// This is a comment",
		},
		{
			name:        "indented",
			indentLevel: 1,
			node:        "This is a comment\nThis is another comment",
			res: stringsutil.StripMargin(`
				|/// This is a comment
				|  /// This is another comment
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			ctx := ctxWithIndentLevel(context.Background(), test.indentLevel)

			// When
			result, err := test.node.Marshal(ctx)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, test.res, string(result))
		})
	}
}

func TestShebangCommentMarshal(t *testing.T) {
	node := ShebangComment("/bin/balsamic eval")
	result, err := node.Marshal(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, `#! /bin/balsamic eval`, string(result))
}
