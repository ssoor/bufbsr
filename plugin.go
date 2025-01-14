package codegenerator

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("marshaling plugin request: %w", err)
	}

	stdout := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, p.Path)

	cmd.Stdin = bytes.NewReader(in)
	cmd.Stderr = io.Discard
	cmd.Stdout = stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	res := &pluginpb.CodeGeneratorResponse{}
	if err := proto.Unmarshal(stdout.Bytes(), res); err != nil {
		return nil, fmt.Errorf("unmarshaling plugin response: %w", err)
	}

	return res, nil
}
