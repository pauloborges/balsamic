package ast

import (
	"context"
	"testing"

	"github.com/pauloborges/balsamic/internal/stringsutil"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationMarshal(t *testing.T) {
	tests := []struct {
		name string
		node *Annotation
		res  string
		err  error
	}{
		{
			name: "annotation",
			node: &Annotation{
				Name: QualifiedIdentifier("test.Annotation"),
			},
			res: `@test.Annotation`,
		},
		{
			name: "with body",
			node: &Annotation{
				Name: QualifiedIdentifier("test.Annotation"),
				Body: &ObjectBody{
					Members: []ObjectMember{
						&ObjectProperty{
							Name:  "key",
							Value: StringExpression("value"),
						},
					},
				},
			},
			res: stringsutil.StripMargin(`
				|@test.Annotation {
				|  key = "value"
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
