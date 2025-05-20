package pkl

import (
	"fmt"
	"io/fs"

	"github.com/pauloborges/balsamic/ast"
	"github.com/pauloborges/balsamic/internal/fsutil/memfs"
)

type Project struct {
	Name    string
	Modules map[string]*Module
}

func NewProject(name string) *Project {
	manifest := &Module{
		Path: "PklProject",
		AST: &ast.Module{
			ParentRelationship: ast.ModuleRelationshipAmends,
			ParentName:         "pkl:Project",
		},
	}

	return &Project{
		Name: name,
		Modules: map[string]*Module{
			manifest.Path: manifest,
		},
	}
}

func (p *Project) AddModule(m *Module) {
	p.Modules[m.Path] = m
}

func (p *Project) Render() (fs.FS, error) {
	fsys := memfs.New()

	for _, m := range p.Modules {
		data, err := m.Marshal()
		if err != nil {
			return nil, fmt.Errorf("marshal Pkl module %s: %w", m.Path, err)
		}

		err = fsys.WriteFile(m.Path, data, 0666)
		if err != nil {
			return nil, fmt.Errorf("write Pkl module %s: %w", m.Path, err)
		}
	}

	return fsys, nil
}
