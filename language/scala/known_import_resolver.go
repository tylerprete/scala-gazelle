package scala

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"

	"github.com/stackb/scala-gazelle/pkg/resolver"
)

// newKnownImportResolver constructs the top-level known import resolver.
func newKnownImportResolver(known resolver.KnownImportRegistry) resolver.KnownImportResolver {
	chain := resolver.NewChainResolver(
		// override resolver is the least performant!
		resolver.NewMemoResolver(resolver.NewOverrideResolver(scalaLangName)),
		resolver.NewKnownResolver(known),
		resolver.NewCrossResolver(scalaLangName),
	)
	scala := resolver.NewScalaResolver(chain)
	return resolver.NewMemoResolver(scala)
}

// ResolveKnownImport implements the resolver.ImportResolver interface.
func (sl *scalaLang) ResolveKnownImport(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, imp string) (*resolver.KnownImport, error) {
	return sl.knownImportResolver.ResolveKnownImport(c, ix, from, lang, imp)
}
