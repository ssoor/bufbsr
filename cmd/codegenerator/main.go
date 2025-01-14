package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CGA1123/codegenerator"
	"github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	path, ok := os.LookupEnv("CODEGENERATOR_REGISTRY_PATH")
	if !ok || path == "" {
		log.Fatalf("CODEGENERATOR_REGISTRY_PATH is not set")
	}

	registry := codegenerator.LocalRegistry(path)
	service := &codegenerator.Service{Registry: registry}

	mux := http.NewServeMux()
	path, handler := registryv1alpha1connect.NewCodeGenerationServiceHandler(service)
	mux.Handle(path, handler)

	// TODO: autogenerate these certs for development.
	// In production, run non-TLS server.
	// LB handles it.
	err := http.ListenAndServeTLS(
		"localhost:1123",
		".local/certstrap/codegenerator.crt",
		".local/certstrap/codegenerator.key",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
