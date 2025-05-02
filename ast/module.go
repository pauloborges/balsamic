package ast

import (
	"context"
	"strconv"

	"github.com/pauloborges/balsamic/internal/bytesutil"
)

type ModuleRelationship string

const (
	ModuleRelationshipExtends ModuleRelationship = "extends"
	ModuleRelationshipAmends  ModuleRelationship = "amends"
)

type Module struct {
	// Optional.
	ShebangComment ShebangComment
	// Optional.
	Docs Docs
	// Optional.
	Annotations Annotations
	// Optional.
	Modifiers Modifiers
	// Required if ParentName/ParentRelationship isn't set.
	Name QualifiedIdentifier
	// Required if Name isn't set, optional otherwise.
	ParentRelationship ModuleRelationship
	// Required if Name isn't set and ParentRelationship is set,
	// optional otherwise.
	ParentName string
	// Optional.
	Imports ImportClauses
	// Optional.
	Members ModuleMembers
}

func (m *Module) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	shebang, err := m.ShebangComment.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(shebang, "\n")

	docs, err := m.Docs.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(docs, "\n")

	annotations, err := m.Annotations.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(annotations, "\n")

	modifiers, err := m.Modifiers.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithSuffix(modifiers, " ")

	if m.Name != "" {
		name, err := m.Name.Marshal(ctx)
		if err != nil {
			return nil, err
		}
		b.WriteWithPrefix("module ", name)
	}

	if m.ParentName != "" {
		if m.Name != "" {
			b.WriteRune(' ')
		}
		b.WriteString(string(m.ParentRelationship))
		b.WriteRune(' ')
		b.WriteString(strconv.Quote(m.ParentName))
	}

	b.WriteRune('\n')

	imports, err := m.Imports.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("\n", imports, "\n")

	members, err := m.Members.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	b.WriteWithPrefixSuffix("\n", members, "\n")

	return b.Bytes(), nil
}

type ImportClause struct {
	// Required.
	Path string
	// Optional.
	Alias string
	// Optional.
	Glob bool
}

func (i *ImportClause) Marshal(ctx context.Context) ([]byte, error) {
	var b bytesutil.Buffer

	if i.Glob {
		b.WriteString("import* ")
	} else {
		b.WriteString("import ")
	}

	b.WriteString(strconv.Quote(i.Path))

	if i.Alias != "" {
		b.WriteString(" as ")
		b.WriteString(i.Alias)
	}

	return b.Bytes(), nil
}

type ImportClauses []*ImportClause

func (i ImportClauses) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, i, "\n")
}

type ModuleMember interface {
	Node
	isModuleMember()
}

type ModuleMembers []ModuleMember

func (m ModuleMembers) Marshal(ctx context.Context) ([]byte, error) {
	return joinNodes(ctx, m, "\n\n")
}
