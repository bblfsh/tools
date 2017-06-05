package main

import "github.com/bblfsh/tools"

type Npath struct {
	Common
}

func (c *Npath) Execute(args []string) error {
	return c.execute(args, tools.Npath{})
}
