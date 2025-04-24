package docker

import (
	"fmt"
	"path/filepath"

	v1alpha1 "github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1"
	"github.com/CGA1123/codegenerator/plugin"
	"github.com/CGA1123/codegenerator/plugin/local"
)

// LocalRegistry reads the available plugins from the folder structure at
// "path".
//
// It will panic if the folder structure does not match our expectations.
//
// The expectation is that for a remote plugin request of
// `<host>/<owner>/<plugin>:<version>`
//
// e.g.
// ```
// # buf.gen.yaml
// plugins:
// - remote: <host>/<owner>/<plugin>:<version>
// ```
//
// There is an executable file at `<owner>/<plugin>/<version>/<plugin>`
//
// <version> is required to match `v1.2.3` (or `/v\d+\.\d+\.\d+`).
func DockerRegistry(path string) *Registry {
	return &Registry{registry: path}
}

// Registry is the container which points to all available plugins.
type Registry struct {
	registry string
}

// Get gets a plugin, if registered.
//
// * Version must be set.
// * Revision must not be set.
func (r *Registry) Get(ref *v1alpha1.CuratedPluginReference) (plugin.Plugin, error) {
	if ref.GetRevision() != 0 {
		return nil, fmt.Errorf("setting version revision is not supported: got revision %v", ref.GetRevision())
	}

	if ref.GetVersion() == "" {
		return nil, fmt.Errorf("not setting a version is not supported")
	}

	pluginRef := fmt.Sprintf("plugins-%s-%s:%s", ref.GetOwner(), ref.GetName(), ref.GetVersion())

	p := &local.Plugin{
		Path:    "docker",
		Args:    []string{"run", "--rm", "-i", filepath.Join(r.registry, pluginRef)},
		Name:    ref.GetName(),
		Version: ref.GetVersion(),
	}

	return p, nil
}
