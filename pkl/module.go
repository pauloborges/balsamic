package pkl

import (
	"context"

	"github.com/pauloborges/balsamic/ast"
)

type Module struct {
	Path string
	AST  *ast.Module
}

func (m *Module) Marshal() ([]byte, error) {
	return m.AST.Marshal(context.Background())
}
