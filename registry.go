package codegenerator

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	v1alpha1 "github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1"
)

func LocalRegistry(path string) *Registry {
	r, err := buildLocalRegistry(path)
	if err != nil {
		log.Fatalf("building local registry: %v", err)
	}

	return r
}

var semverRegex = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func buildLocalRegistry(path string) (*Registry, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("expanding path: %w", err)
	}

	slog.Info("building local registry", "path", path)

	registry := map[string]map[string]map[string]*Plugin{}

	owners, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("listing registry path: %w", err)
	}

	for _, owner := range owners {
		ownerName := owner.Name()

		if !owner.IsDir() {
			return nil, fmt.Errorf("expected %s/%s to be a directory", path, ownerName)
		}

		plugins, err := os.ReadDir(filepath.Join(path, ownerName))
		if err != nil {
			return nil, fmt.Errorf("listing plugins for %s: %w", ownerName, err)
		}

		for _, plugin := range plugins {
			pluginName := plugin.Name()

			if !plugin.IsDir() {
				return nil, fmt.Errorf("expected %s/%s/%s to be a directory", path, ownerName, pluginName)
			}

			versions, err := os.ReadDir(
				filepath.Join(path, ownerName, pluginName),
			)
			if err != nil {
				return nil, fmt.Errorf("reading versions for %s/%s: %w", ownerName, pluginName, err)
			}

			for _, version := range versions {
				versionName := version.Name()

				if !semverRegex.MatchString(versionName) {
					return nil, fmt.Errorf("incorrect version path: %s", filepath.Join(pluginName, versionName))
				}

				if !version.IsDir() {
					return nil, fmt.Errorf("expected %s/%s/%s/%s to be a directory", path, ownerName, pluginName, versionName)
				}

				info, err := os.Stat(filepath.Join(
					path,
					ownerName,
					pluginName,
					versionName,
					pluginName,
				))
				if err != nil {
					return nil, fmt.Errorf("stating binary for %s/%s@%s: %w", ownerName, pluginName, versionName, err)
				}

				slog.Info("found plugin", "owner", ownerName, "plugin", pluginName, "version", versionName, "name", fmt.Sprintf("%s/%s:%s", ownerName, pluginName, versionName))

				execable := (info.Mode().Perm() & 0111) != 0
				if !execable {
					return nil, fmt.Errorf("plugin %s/%s@%s (%s) is not executable", ownerName, pluginName, versionName, filepath.Join(
						path, pluginName, versionName, pluginName,
					))
				}

				p := &Plugin{
					Path:    filepath.Join(path, ownerName, pluginName, versionName, pluginName),
					Name:    pluginName,
					Version: versionName,
				}

				if _, ok := registry[ownerName]; !ok {
					registry[ownerName] = map[string]map[string]*Plugin{}
				}

				if _, ok := registry[pluginName]; !ok {
					registry[ownerName][pluginName] = map[string]*Plugin{}
				}

				registry[ownerName][pluginName][versionName] = p
			}
		}
	}

	return &Registry{registry: registry}, nil
}

// Registry is the container which points to all available plugins.
type Registry struct {
	registry map[string]map[string]map[string]*Plugin
}

func (r *Registry) Get(ref *v1alpha1.CuratedPluginReference) (*Plugin, error) {
	if ref.GetRevision() != 0 {
		return nil, fmt.Errorf("setting version revision is not supported: got revision %v", ref.GetRevision())
	}

	pluginRef := fmt.Sprintf("%s/%s:%s", ref.GetOwner(), ref.GetName(), ref.GetVersion())

	plugins, ok := r.registry[ref.GetOwner()]
	if !ok {
		return nil, fmt.Errorf("plugin not found '%s': owner not found", pluginRef)
	}

	versions, ok := plugins[ref.GetName()]
	if !ok {
		return nil, fmt.Errorf("plugin not found '%s': plugin not found", pluginRef)
	}

	plugin, ok := versions[ref.GetVersion()]
	if !ok {
		return nil, fmt.Errorf("plugin not found '%s': version not found", pluginRef)
	}

	return plugin, nil
}
