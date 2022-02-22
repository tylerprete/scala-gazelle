package scala

import (
	"errors"
	"flag"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/resolve"
)

// ConfigurableCrossResolver implementations support the CrossResolver interface
// as well as a subset of config.Configurer.  This interface is provided to
// support different implementations of a scala cross-resolver.  A simple
// implementation might be based off a CSV file, whereas a larger monorepo may
// require a more sophisticated cache.
type ConfigurableCrossResolver interface {
	resolve.CrossResolver
	// RegisterFlags implements part of the config.Configurer interface.
	RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config)
	// CheckFlags implements part of the config.Configurer interface.
	CheckFlags(fs *flag.FlagSet, c *config.Config) error
}

// GazellePhaseTransitionListener is an optional interface for a cross-resolver
// that wants phase transition notification.  Errors are considered fatal.
type GazellePhaseTransitionListener interface {
	OnIndexPhase() error
	OnResolvePhase() error
}

// ErrUnknownResolver is the error returned when a CrossResolver is not known.
var ErrUnknownResolver = errors.New("unknown CrossResolver")

// CrossResolverRegistry represents a mapping of Cross Resolver implementations.
type CrossResolverRegistry interface {
	// CrossResolverNames returns a sorted list of CrossResolver names.
	CrossResolverNames() []string
	// LookupCrossResolver returns the implementation under the given name.  If the CrossResolver
	// is not found, ErrUnknownResolver is returned.
	LookupCrossResolver(name string) (ConfigurableCrossResolver, error)
	// MustRegisterCrossResolver installs a ConfigurableCrossResolver implementation under the given
	// name in the global CrossResolver registry.  Panic will occur if the same CrossResolver is
	// registered multiple times.
	MustRegisterCrossResolver(name string, resolver ConfigurableCrossResolver) CrossResolverRegistry
}

// CrossResolvers returns a reference to the global CrossResolverRegistry
func CrossResolvers() CrossResolverRegistry {
	return globalCrossResolverRegistry
}

// registry is the default registry singleton.
var globalCrossResolverRegistry = &crossResolverRegistry{
	crossResolvers: make(map[string]ConfigurableCrossResolver),
}

// crossResolverRegistry implements CrossResolverRegistry.
type crossResolverRegistry struct {
	crossResolvers map[string]ConfigurableCrossResolver
}

// CrossResolverNames implements part of the CrossResolverRegistry interface.
func (p *crossResolverRegistry) CrossResolverNames() []string {
	names := make([]string, 0)
	for name := range p.crossResolvers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MustRegisterCrossResolver implements part of the ResolverRegistry interface.
func (p *crossResolverRegistry) MustRegisterCrossResolver(name string, resolver ConfigurableCrossResolver) CrossResolverRegistry {
	_, ok := p.crossResolvers[name]
	if ok {
		panic("duplicate CrossResolver registration: " + name)
	}
	p.crossResolvers[name] = resolver
	return p
}

// LookupCrossResolver implements part of the ResolverRegistry interface.
func (p *crossResolverRegistry) LookupCrossResolver(name string) (ConfigurableCrossResolver, error) {
	resolver, ok := p.crossResolvers[name]
	if !ok {
		return nil, ErrUnknownResolver
	}
	return resolver, nil
}
