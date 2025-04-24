package registry

import (
	v1alpha1 "github.com/CGA1123/codegenerator/gen/buf/alpha/registry/v1alpha1"
	"github.com/CGA1123/codegenerator/plugin"
)

type Registry interface {
	Get(ref *v1alpha1.CuratedPluginReference) (plugin.Plugin, error)
}
