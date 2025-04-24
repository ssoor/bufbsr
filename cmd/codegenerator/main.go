package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/CGA1123/codegenerator"
	"github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/CGA1123/codegenerator/registry"
	"github.com/CGA1123/codegenerator/registry/docker"
	"github.com/CGA1123/codegenerator/registry/local"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	var (
		typ     = flag.String("type", "docker", "The types of the registry support docker and local")
		address = flag.String("address", "0.0.0.0:443", "The address listened for by the service")
		tlsCrt  = flag.String("tls-crt", ".local/certstrap/codegenerator.crt", "The certificate used by TLS")
		tlsKey  = flag.String("tls-key", ".local/certstrap/codegenerator.key", "The certificate private key used by TLS")
	)
	flag.Parse()

	path, _ := os.LookupEnv("CODEGENERATOR_REGISTRY_PATH")

	var registry registry.Registry
	switch *typ {
	case "local":
		if path == "" {
			log.Fatalf("CODEGENERATOR_REGISTRY_PATH is not set")
		}
		registry = local.LocalRegistry(path)
	default:
		registry = docker.DockerRegistry(path)
	}
	service := &codegenerator.Service{Registry: registry}

	mux := http.NewServeMux()
	path, handler := registryv1alpha1connect.NewCodeGenerationServiceHandler(service)
	mux.Handle(path, handler)

	log.Println("server listen address:", *address)
	ln, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("listen address err: %v", err)
	}

	defer ln.Close()

	if *tlsCrt != "" && *tlsKey != "" {
		err = http.ServeTLS(
			ln,
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
			*tlsCrt,
			*tlsKey,
		)
	} else {
		err = http.Serve(
			ln,
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		)
	}

	if err != nil {
		log.Fatalf("server running err: %v", err)
	}
}
