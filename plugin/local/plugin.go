package local

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin wraps a plugin binary for local execution.
type Plugin struct {
	Cwd     string
	Path    string
	Args    []string
	Name    string
	Version string
}

func (p *Plugin) Generate(ctx context.Context, req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	in, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling plugin request: %w", err)
	}

	stdout := &bytes.Buffer{}
	errout := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, p.Path, p.Args...)

	cmd.Stdin = bytes.NewReader(in)
	cmd.Stderr = errout // io.Discard
	cmd.Stdout = stdout
	cmd.Dir = p.Cwd

	if err := cmd.Run(); err != nil {
		fmt.Println("execute plugin failed, errout:", errout)
		return nil, err
	}

	res := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(stdout.Bytes(), res); err != nil {
		return nil, fmt.Errorf("unmarshaling plugin response: %w", err)
	}

	return res, nil
}
