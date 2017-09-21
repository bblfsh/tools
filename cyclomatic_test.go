package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func TestCyclomaticComplexity(t *testing.T) {
	require := require.New(t)
	n := &uast.Node{InternalType: "module",
		Children: []*uast.Node{
			{InternalType: "root"}, // 1 (initial)
			// Prefix is the default so it doesnt need any role
			{InternalType: "if1", Roles: []uast.Role{uast.If}, Children: []*uast.Node{ // 2 (If)
				{InternalType: "if1else1", Roles: []uast.Role{uast.IfElse}, Children: []*uast.Node{ // 0
					{InternalType: "if1else1foreach", Roles: []uast.Role{uast.ForEach}, Children: []*uast.Node{ // 3 (ForEach)
						{InternalType: "foreach_child1"},                                             // 0
						{InternalType: "foreach_child2_continue", Roles: []uast.Role{uast.Continue}}, // 4 (Continue)
					}},
					{InternalType: "if1else1if", Roles: []uast.Role{uast.If}, Children: []*uast.Node{ // 5 (If)
						{InternalType: "elseif_child1"},                                // 0
						{InternalType: "opAnd", Roles: []uast.Role{uast.OpBooleanAnd}}, // 6 (OpBooleanAnd)
						{InternalType: "elseif_child2"},                                // 0
					}},
				}},
				{InternalType: "break", Roles: []uast.Role{uast.Break}},
			},
			}}}
	require.Equal(cyclomaticComplexity(n), 6)
}
