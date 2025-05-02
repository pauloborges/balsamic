package ast

import (
	"bytes"
	"context"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type Comment interface {
	Node
	isComment()
}

type LineComment string

func (c LineComment) isComment() {}

func (c LineComment) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer
	prefix := indentation(ctx) + "// "
	b.WriteWithPrefix(prefix, []byte(c))
	return b.Bytes(), nil
}

type BlockComment string

func (c BlockComment) isComment() {}

func (c BlockComment) Marshal(ctx context.Context) ([]byte, error) {
	var b bytes.Buffer

	if c == "" {
		return nil, nil
	}

	b.WriteString("/*")
	b.WriteString(newlineWithIndentation(ctx))

	for i, line := range bytes.Split([]byte(string(c)), []byte{'\n'}) {
		if i > 0 {
			b.WriteString(newlineWithIndentation(ctx))
		}
		if len(line) > 0 {
			b.WriteString("  ")
		}
		b.Write(line)
	}

	b.WriteString(newlineWithIndentation(ctx))
	b.WriteString("*/")

	return b.Bytes(), nil
}

type Docs string

func (c Docs) isComment() {}

func (c Docs) Marshal(ctx context.Context) ([]byte, error) {
	if c == "" {
		return nil, nil
	}

	var b bytes.Buffer

	for i, line := range bytes.Split([]byte(string(c)), []byte{'\n'}) {
		if i > 0 {
			b.WriteString(newlineWithIndentation(ctx))
		}
		if len(line) > 0 {
			b.WriteString("/// ")
			b.Write(line)
		} else {
			b.WriteString("///")
		}
	}

	return b.Bytes(), nil
}

type ShebangComment string

func (c ShebangComment) isComment() {}

func (c ShebangComment) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer
	b.WriteWithPrefix("#! ", []byte(c))
	return b.Bytes(), nil
}
