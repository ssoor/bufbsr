package plugin

import (
	"context"

	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin interface {
	Generate(ctx context.Context, req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}
