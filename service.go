package codegenerator

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	imagev1 "github.com/CGA1123/codegenerator/gen/buf/alpha/image/v1"
	v1alpha1 "github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1"
	"github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/bufbuild/protoplugin/protopluginutil"
)

var _ registryv1alpha1connect.CodeGenerationServiceHandler = (*Service)(nil)

type Service struct {
	Registry *Registry
}

func (s *Service) GenerateCode(
	ctx context.Context,
	req *connect.Request[v1alpha1.GenerateCodeRequest],
) (*connect.Response[v1alpha1.GenerateCodeResponse], error) {
	msg := req.Msg
	responses := make([]*v1alpha1.PluginGenerationResponse, len(msg.GetRequests()))
	for i, pluginRequest := range msg.GetRequests() {
		plugin, err := s.Registry.Get(pluginRequest.GetPluginReference())
		if err != nil {
			return nil, err
		}

		genReq, err := imageToCodeGenerationRequest(msg.GetImage(), pluginRequest)
		if err != nil {
			return nil, err
		}

		pluginResponse, err := plugin.Generate(ctx, genReq)
		if err != nil {
			return nil, err
		}

		responses[i] = &v1alpha1.PluginGenerationResponse{Response: pluginResponse}
	}

	return connect.NewResponse(
		&v1alpha1.GenerateCodeResponse{
			Responses: responses,
		}), nil
}

func imageToCodeGenerationRequest(image *imagev1.Image, plug *v1alpha1.PluginGenerationRequest) (*pluginpb.CodeGeneratorRequest, error) {
	imageFiles := image.GetFile()
	parameter := plug.GetOptions()
	request := &pluginpb.CodeGeneratorRequest{
		ProtoFile: make([]*descriptorpb.FileDescriptorProto, len(imageFiles)),
		Parameter: proto.String(strings.Join(parameter, ",")),
	}

	for i, imageFile := range imageFiles {
		fileDescriptorProto := ImageFileToDescriptor(imageFile)

		// TODO: check if we should generate this file?
		request.FileToGenerate = append(request.FileToGenerate, imageFile.GetName())
		request.SourceFileDescriptors = append(request.SourceFileDescriptors, fileDescriptorProto)
		var err error
		fileDescriptorProto, err = protopluginutil.StripSourceRetentionOptions(fileDescriptorProto)
		if err != nil {
			return nil, fmt.Errorf("failed to strip source-retention options for file %q when constructing a CodeGeneratorRequest: %w", imageFile.GetName(), err)
		}

		request.ProtoFile[i] = fileDescriptorProto
	}

	return request, nil
}

func ImageFileToDescriptor(img *imagev1.ImageFile) *descriptorpb.FileDescriptorProto {
	return &descriptorpb.FileDescriptorProto{
		Name:             img.Name,
		Package:          img.Package,
		Dependency:       img.Dependency,
		PublicDependency: img.PublicDependency,
		WeakDependency:   img.WeakDependency,
		MessageType:      img.MessageType,
		EnumType:         img.EnumType,
		Service:          img.Service,
		Extension:        img.Extension,
		Options:          img.Options,
		SourceCodeInfo:   img.SourceCodeInfo,
		Syntax:           img.Syntax,
		Edition:          img.Edition,
	}
}
