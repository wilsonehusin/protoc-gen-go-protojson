package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		return fmt.Errorf("decoding request: %w", err)
	}
	resp, err := generate(req)
	if err != nil {
		return fmt.Errorf("generating: %w", err)
	}
	out, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("encoding response: %w", err)
	}
	if _, err = os.Stdout.Write(out); err != nil {
		return fmt.Errorf("writing to stdout: %w", err)
	}
	return nil
}

//go:embed protojson.go.tmpl
var rawTemplate string

var tmpl = template.Must(template.New("protojson").Parse(rawTemplate))

type Values struct {
	GoPackageName string
	GoPackagePath string
	ProtoMessages []ProtoMessage
}

type ProtoMessage struct {
	Name string
}

func generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	resp := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)),
		MinimumEdition:    proto.Int32(int32(descriptorpb.Edition_EDITION_PROTO2)),
		MaximumEdition:    proto.Int32(int32(descriptorpb.Edition_EDITION_2024)),
	}

	for _, protofile := range req.GetProtoFile() {
		var buf bytes.Buffer
		pkg := protofile.GetOptions().GetGoPackage()
		path, base, _ := strings.Cut(pkg, ";")
		if base == "" {
			base = filepath.Base(path)
		}
		v := &Values{
			GoPackageName: base,
			GoPackagePath: path,
			ProtoMessages: []ProtoMessage{},
		}

		for _, msg := range protofile.GetMessageType() {
			v.ProtoMessages = append(v.ProtoMessages, ProtoMessage{
				Name: msg.GetName(),
			})
		}

		if err := tmpl.Execute(&buf, v); err != nil {
			return nil, fmt.Errorf("executing template: %w", err)
		}

		name, _, _ := strings.Cut(protofile.GetName(), ".")
		name += ".protojson.go"

		// Return a file to protoc to write.
		file := &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(buf.String()),
		}
		resp.File = append(resp.File, file)
	}
	return resp, nil
}
