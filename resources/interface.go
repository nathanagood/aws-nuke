package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type ResourceListers map[string]ResourceLister

type ResourceLister func(s *session.Session) ([]Resource, error)

type Resource interface {
	Remove() error
}

type Filter interface {
	Resource
	Filter() error
}

type LegacyStringer interface {
	Resource
	String() string
}

type ResourcePropertyGetter interface {
	Resource
	Properties() types.Properties
}

type FeatureFlagGetter interface {
	Resource
	FeatureFlags(config.FeatureFlags)
}

var resourceListers = make(ResourceListers)

// [Option 1]: Add a "deps" parameter here to register dependencies.
// This isn't too messy, as the developer writing the code to handle
// the resource can add this at the time. It can default to nothing,
// so backwards compatibility is perserved.
func register(name string, lister ResourceLister, deps ...string) {
	_, exists := resourceListers[name]
	if exists {
		panic(fmt.Sprintf("a resource with the name %s already exists", name))
	}

	// [Option 1]: Change this to a list of structures?
	resourceListers[name] = lister
}

func GetListers() ResourceListers {
	return resourceListers
}

func GetLister(name string) ResourceLister {
	return resourceListers[name]
}

func GetListerNames() []string {
	names := []string{}
	for resourceType, _ := range GetListers() {
		names = append(names, resourceType)
	}

	return names
}
