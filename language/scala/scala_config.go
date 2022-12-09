package scala

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"

	"github.com/stackb/scala-gazelle/pkg/collections"
	"github.com/stackb/scala-gazelle/pkg/resolver"
	"github.com/stackb/scala-gazelle/pkg/scalarule"
)

type annotation int

const (
	AnnotateUnknown        annotation = 0
	AnnotateImports        annotation = 1
	AnnotateUnresolvedDeps annotation = 2
	AnnotateResolvedDeps   annotation = 3
)

const (
	// scalaAnnotate is the name of a directive.
	scalaAnnotateDirective = "scala_annotate"
	// ruleDirective is the directive for enabling rule generation.
	ruleDirective = "scala_rule"
	// resolveGlobDirective implements override via globs.
	resolveGlobDirective = "resolve_glob"
	// resolveWithDirective adds additional imports for resolution
	resolveWithDirective = "resolve_with"
	// resolveKindRewriteName allows renaming of resolved labels.
	resolveKindRewriteName = "resolve_kind_rewrite_name"
)

// scalaConfig represents the config extension for the a scala package.
type scalaConfig struct {
	config            *config.Config
	rel               string
	universe          resolver.Universe
	overrides         []*overrideSpec
	implicitImports   []*implicitImportSpec
	rules             map[string]*scalarule.Config
	labelNameRewrites map[string]resolver.LabelNameRewriteSpec
	annotations       map[annotation]interface{}
}

// newScalaConfig initializes a new scalaConfig.
func newScalaConfig(universe resolver.Universe, config *config.Config, rel string) *scalaConfig {
	return &scalaConfig{
		config:            config,
		rel:               rel,
		universe:          universe,
		annotations:       make(map[annotation]interface{}),
		labelNameRewrites: make(map[string]resolver.LabelNameRewriteSpec),
		rules:             make(map[string]*scalarule.Config),
	}
}

// getScalaConfig returns the scala config.  Can be nil.
func getScalaConfig(config *config.Config) *scalaConfig {
	if existingExt, ok := config.Exts[scalaLangName]; ok {
		return existingExt.(*scalaConfig)
	} else {
		return nil
	}
}

// getOrCreateScalaConfig either inserts a new config into the map under the
// language name or replaces it with a clone.
func getOrCreateScalaConfig(universe resolver.Universe, config *config.Config, rel string) *scalaConfig {
	var cfg *scalaConfig
	if existingExt, ok := config.Exts[scalaLangName]; ok {
		cfg = existingExt.(*scalaConfig).clone(config, rel)
	} else {
		cfg = newScalaConfig(universe, config, rel)
	}
	config.Exts[scalaLangName] = cfg
	return cfg
}

// clone copies this config to a new one.
func (c *scalaConfig) clone(config *config.Config, rel string) *scalaConfig {
	clone := newScalaConfig(c.universe, config, rel)
	for k, v := range c.annotations {
		clone.annotations[k] = v
	}
	for k, v := range c.rules {
		clone.rules[k] = v.Clone()
	}
	for k, v := range c.labelNameRewrites {
		clone.labelNameRewrites[k] = v
	}
	if c.overrides != nil {
		clone.overrides = c.overrides[:]
	}
	if c.implicitImports != nil {
		clone.implicitImports = c.implicitImports[:]
	}
	return clone
}

func (c *scalaConfig) canProvide(from label.Label) bool {
	for _, provider := range c.universe.SymbolProviders() {
		if provider.CanProvide(from, c.universe.GetKnownRule) {
			return true
		}
	}
	return false
}

// GetKnownRule translates relative labels into their absolute form.
func (c *scalaConfig) GetKnownRule(from label.Label) (*rule.Rule, bool) {
	if from.Name == "" {
		return nil, false
	}
	if from.Repo == "" {
		from = label.New(c.config.RepoName, from.Pkg, from.Name)
	}
	if from.Pkg == "" && from.Repo == c.config.RepoName {
		from = label.New(from.Repo, c.rel, from.Name)
	}
	return c.universe.GetKnownRule(from)
}

// parseDirectives is called in each directory visited by gazelle.  The relative
// directory name is given by 'rel' and the list of directives in the BUILD file
// are specified by 'directives'.
func (c *scalaConfig) parseDirectives(directives []rule.Directive) (err error) {
	for _, d := range directives {
		switch d.Key {
		case ruleDirective:
			err = c.parseRuleDirective(d)
			if err != nil {
				return fmt.Errorf(`invalid directive: "gazelle:%s %s": %w`, d.Key, d.Value, err)
			}
		case resolveGlobDirective:
			c.parseResolveGlobDirective(d)
		case resolveWithDirective:
			c.parseResolveWithDirective(d)
		case resolveKindRewriteName:
			c.parseResolveKindRewriteNameDirective(d)
		case scalaAnnotateDirective:
			if err := c.parseScalaAnnotation(d); err != nil {
				return err
			}
		}
	}
	return
}

func (c *scalaConfig) parseRuleDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) < 3 {
		return fmt.Errorf("expected three or more fields, got %d", len(fields))
	}
	name, param, value := fields[0], fields[1], strings.Join(fields[2:], " ")
	r, err := c.getOrCreateScalaRuleConfig(c.config, name)
	if err != nil {
		return err
	}
	return r.ParseDirective(name, param, value)
}

func (c *scalaConfig) parseResolveGlobDirective(d rule.Directive) {
	parts := strings.Fields(d.Value)
	o := overrideSpec{}
	var lbl string
	if len(parts) != 4 {
		return
	}
	if parts[0] != scalaLangName {
		return
	}

	o.imp.Lang = parts[0]
	o.lang = parts[1]
	o.imp.Imp = parts[2]
	lbl = parts[3]

	var err error
	o.dep, err = label.Parse(lbl)
	if err != nil {
		log.Fatalf("bad gazelle:%s directive value %q: %v", resolveGlobDirective, d.Value, err)
		return
	}
	c.overrides = append(c.overrides, &o)
}

func (c *scalaConfig) parseResolveWithDirective(d rule.Directive) {
	parts := strings.Fields(d.Value)
	if len(parts) < 3 {
		log.Printf("invalid gazelle:%s directive: expected 3+ parts, got %d (%v)", resolveWithDirective, len(parts), parts)
		return
	}
	c.implicitImports = append(c.implicitImports, &implicitImportSpec{
		lang: parts[0],
		imp:  parts[1],
		deps: parts[2:],
	})
}

func (c *scalaConfig) parseResolveKindRewriteNameDirective(d rule.Directive) {
	parts := strings.Fields(d.Value)
	if len(parts) != 3 {
		log.Printf("invalid gazelle:%s directive: expected [KIND SRC_NAME DST_NAME], got %v", resolveKindRewriteName, parts)
		return
	}
	kind := parts[0]
	src := parts[1]
	dst := parts[2]

	c.labelNameRewrites[kind] = resolver.LabelNameRewriteSpec{Src: src, Dst: dst}
}

func (c *scalaConfig) parseScalaAnnotation(d rule.Directive) error {
	for _, key := range strings.Fields(d.Value) {
		intent := collections.ParseIntent(key)
		annot := parseAnnotation(intent.Value)
		if annot == AnnotateUnknown {
			return fmt.Errorf("invalid directive gazelle:%s: unknown annotation value '%v'", d.Key, intent.Value)
		}
		if intent.Want {
			var val interface{}
			c.annotations[annot] = val
		} else {
			delete(c.annotations, annot)
		}
	}
	return nil
}

func (c *scalaConfig) getOrCreateScalaRuleConfig(config *config.Config, name string) (*scalarule.Config, error) {
	r, ok := c.rules[name]
	if !ok {
		r = scalarule.NewConfig(config, name)
		r.Implementation = name
		c.rules[name] = r
	}
	return r, nil
}

func (c *scalaConfig) getImplicitImports(lang, imp string) (deps []string) {
	for _, d := range c.implicitImports {
		if d.lang != lang {
			continue
		}
		if d.imp != imp {
			continue
		}
		deps = append(deps, d.deps...)
	}
	return
}

// configuredRules returns an ordered list of configured rules
func (c *scalaConfig) configuredRules() []*scalarule.Config {

	names := make([]string, 0)
	for name := range c.rules {
		names = append(names, name)
	}
	sort.Strings(names)
	rules := make([]*scalarule.Config, len(names))
	for i, name := range names {
		rules[i] = c.rules[name]
	}
	return rules
}

func (c *scalaConfig) shouldAnnotateImports() bool {
	_, ok := c.annotations[AnnotateImports]
	return ok
}

func (c *scalaConfig) shouldAnnotateResolvedDeps() bool {
	_, ok := c.annotations[AnnotateResolvedDeps]
	return ok
}

func (c *scalaConfig) shouldAnnotateUnresolvedDeps() bool {
	_, ok := c.annotations[AnnotateUnresolvedDeps]
	return ok
}

func (c *scalaConfig) Comment() build.Comment {
	return build.Comment{Token: "# " + c.String()}
}

func (c *scalaConfig) String() string {
	return fmt.Sprintf("scalaConfig rel=%q, annotations=%+v", c.rel, c.annotations)
}

type overrideSpec struct {
	imp  resolve.ImportSpec
	lang string
	dep  label.Label
}

type implicitImportSpec struct {
	// lang is the language to which this implicit applies.  Always 'scala' for now.
	lang string
	// imp is the "source" dependency (e.g. LazyLogging)
	imp string
	// dep is the "destination" dependencies (e.g. org.slf4j.Logger)
	deps []string
}

func parseAnnotation(val string) annotation {
	switch val {
	case "imports":
		return AnnotateImports
	case "resolved_deps":
		return AnnotateResolvedDeps
	case "unresolved_deps":
		return AnnotateUnresolvedDeps
	default:
		return AnnotateUnknown
	}
}

// isSameImport returns true if the "from" and "to" labels are the same,
// normalizing to the config.RepoName and performing label name remapping if the
// kind matches.
func isSameImport(sc *scalaConfig, kind string, from, to label.Label) bool {
	if from.Repo == "" {
		from = label.New(sc.config.RepoName, from.Pkg, from.Name)
	}
	if to.Repo == "" {
		to = label.New(sc.config.RepoName, to.Pkg, to.Name)
	}
	if mapping, ok := sc.labelNameRewrites[kind]; ok {
		from = mapping.Rewrite(from)
	}
	return from == to
}
