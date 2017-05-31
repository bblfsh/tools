package tools

import "github.com/bblfsh/sdk/uast"

type Tool interface {
	Exec(*uast.Node) error
}
