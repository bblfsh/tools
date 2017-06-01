package main

import (
	"context"
	"io/ioutil"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/bblfsh/sdk/protocol"
	"github.com/bblfsh/sdk/uast"
	"github.com/bblfsh/tools"
)

type Common struct {
	Address  string `long:"address" description:"server adress to connect to" default:"localhost:9432"`
	Language string `long:"language" description:"language of the input" default:""`
	Args     struct {
		File string `positional-arg-name:"file" required:"true"`
	} `positional-args:"yes"`
}

func (c *Common) execute(args []string, tool tools.Tooler) error {
	logrus.Debugf("executing command")

	request, err := c.buildRequest()
	if err != nil {
		return err
	}

	uast, err := c.parseRequest(request)
	if err != nil {
		return err
	}

	return tool.Exec(uast)
}

func (c *Common) buildRequest() (*protocol.ParseUASTRequest, error) {
	logrus.Debugf("reading file %s", c.Args.File)
	content, err := ioutil.ReadFile(c.Args.File)
	if err != nil {
		return nil, err
	}

	return &protocol.ParseUASTRequest{Content: string(content)}, nil
}

func (c *Common) parseRequest(request *protocol.ParseUASTRequest) (*uast.Node, error) {
	logrus.Debugf("dialing request at %s", c.Address)
	connection, err := grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := protocol.NewProtocolServiceClient(connection)
	response, err := client.ParseUAST(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	return response.UAST, nil
}
