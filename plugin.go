package codegenerator

import (
	"bytes"
	"context"
	"os/exec"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin wraps a plugin binary for local execution.
type Plugin struct {
	Path    string
	Name    string
	Version string
}

func (p *Plugin) Generate(ctx context.Context, req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	in, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, p.Path)

	cmd.Stdin = bytes.NewReader(in)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	// TODO: handle errors, if it's a *exec.ExitError, read stderr and return it?
	icerr := cmd.Run()
	if icerr != nil {
		return nil, err
	}

	res := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(stdout.Bytes(), res); err != nil {
		return nil, err
	}

	return res, nil
}
