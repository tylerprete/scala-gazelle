package scala

import (
	"flag"
	"log"
	"path/filepath"
	"sort"

	"github.com/stackb/rules_proto/pkg/protoc"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	ScalaLangName = "scala"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	sourceResolver := newScalaSourceIndexResolver()
	scalaCompiler := newScalaCompiler()
	importRegistry := newImportRegistry(sourceResolver, scalaCompiler)

	return &scalaLang{
		ruleRegistry:    globalRuleRegistry,
		scalaFileParser: sourceResolver,
		scalaCompiler:   scalaCompiler,
		packages:        make(map[string]*scalaPackage),
		resolvers: []ConfigurableCrossResolver{
			sourceResolver,
			newScalaClassIndexResolver(),
		},
		importRegistry: importRegistry,
	}
}

// scalaLang implements language.Language.
type scalaLang struct {
	// ruleRegistry is the rule registry implementation.  This holds the rules
	// configured via gazelle directives by the user.
	ruleRegistry RuleRegistry
	// importRegistry instance tracks all known info about imports and rules and
	// is used during import disambiguation.
	importRegistry *importRegistry
	// scalaFileParser is the parser implementation.  This is given to each
	// ScalaPackage during GenerateRules such that rule implementations can use
	// it.
	scalaFileParser ScalaFileParser
	// scalaCompiler is the compiler implementation.  This is passed to the
	// importRegistry for use during import disambiguation.
	scalaCompiler *scalaCompiler
	// packages is map from the config.Rel to *scalaPackage for the
	// workspace-relative packate name.
	packages map[string]*scalaPackage
	// isResolvePhase is a flag that is tracks if at least one Resolve() call
	// has occurred.  It can be used to determine when the rule indexing phase
	// has completed and deps resolution phase has started (and calls
	// onResolvePhase).
	isResolvePhase bool
	// resolvers is a list of cross resolver implementations.  Typically there
	// are two: one to help with third-party code, one to help with first-partt
	// code.
	resolvers []ConfigurableCrossResolver
}

// Name returns the name of the language. This should be a prefix of the kinds
// of rules generated by the language, e.g., "go" for the Go extension since it
// generates "go_library" rules.
func (sl *scalaLang) Name() string { return ScalaLangName }

// OnBegin implements part of the language.Lifecycler interface.
func (sl *scalaLang) OnBegin() {
	log.Println("-- BEGIN ---")
}

// OnIndex implements part of the language.Lifecycler interface.
func (sl *scalaLang) OnIndex() {
	log.Println("-- INDEX PHASE ---")
}

// OnResolve implements part of the language.Lifecycler interface.
func (sl *scalaLang) OnResolve() {
	log.Println("-- RESOLVE PHASE ---")
}

// OnEnd implements part of the language.Lifecycler interface.
func (sl *scalaLang) OnEnd() {
	sl.scalaCompiler.stop()
	log.Println("-- END ---")
}

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (sl *scalaLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	getOrCreateScalaConfig(c) // ignoring return value, only want side-effect

	for _, r := range sl.resolvers {
		r.RegisterFlags(fs, cmd, c)
	}
	sl.scalaCompiler.RegisterFlags(fs, cmd, c)
}

// Configure implements part of the config.Configurer
func (sl *scalaLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	for _, r := range sl.resolvers {
		if err := r.CheckFlags(fs, c); err != nil {
			return err
		}
	}
	if err := sl.scalaCompiler.CheckFlags(fs, c); err != nil {
		return err
	}
	return nil
}

// Configure implements part of the config.Configurer
func (*scalaLang) KnownDirectives() []string {
	return []string{
		ruleDirective,
		overrideDirective,
	}
}

// Loads returns .bzl files and symbols they define. Every rule generated by
// GenerateRules, now or in the past, should be loadable from one of these
// files.
func (sl *scalaLang) Loads() []rule.LoadInfo {
	// Merge symbols
	symbolsByLoadName := make(map[string][]string)

	for _, name := range sl.ruleRegistry.RuleNames() {
		rule, err := sl.ruleRegistry.LookupRule(name)
		if err != nil {
			log.Fatal(err)
		}
		load := rule.LoadInfo()
		symbolsByLoadName[load.Name] = append(symbolsByLoadName[load.Name], load.Symbols...)
	}

	// Ensure names are sorted otherwise order of load statements can be
	// non-deterministic
	keys := make([]string, 0)
	for name := range symbolsByLoadName {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// Build final load list
	loads := make([]rule.LoadInfo, 0)
	for _, name := range keys {
		symbols := symbolsByLoadName[name]
		sort.Strings(symbols)
		loads = append(loads, rule.LoadInfo{
			Name:    name,
			Symbols: symbols,
		})
	}
	return loads
}

// Configure implements part of the config.Configurer
func (sl *scalaLang) Configure(c *config.Config, rel string, f *rule.File) {
	if f == nil {
		return
	}
	if err := getOrCreateScalaConfig(c).ParseDirectives(rel, f.Directives); err != nil {
		log.Fatalf("error while parsing rule directives in package %q: %v", rel, err)
	}
}

// Kinds returns a map of maps rule names (kinds) and information on how to
// match and merge attributes that may be found in rules of those kinds. All
// kinds of rules generated for this language may be found here.
func (sl *scalaLang) Kinds() map[string]rule.KindInfo {
	kinds := make(map[string]rule.KindInfo)

	for _, name := range sl.ruleRegistry.RuleNames() {
		rule, err := sl.ruleRegistry.LookupRule(name)
		if err != nil {
			log.Fatal("Kinds:", err)
		}
		kinds[rule.Name()] = rule.KindInfo()
	}

	return kinds
}

// Fix repairs deprecated usage of language-specific rules in f. This is called
// before the file is indexed. Unless c.ShouldFix is true, fixes that delete or
// rename rules should not be performed.
func (sl *scalaLang) Fix(c *config.Config, f *rule.File) {
}

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested in
// depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a struct to
// avoid breaking implementations in the future when new fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
func (sl *scalaLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	if debug {
		log.Println("visiting", args.Rel)
	}

	cfg := getOrCreateScalaConfig(args.Config)

	pkg := newScalaPackage(sl.ruleRegistry, sl.scalaFileParser, sl.importRegistry, args.Rel, args.File, cfg)
	sl.packages[args.Rel] = pkg

	for _, r := range args.OtherGen {
		if r.Kind() != "proto_library2" {
			continue
		}
		if !hasPackageProto(args.RegularFiles) {
			continue
		}
		srcs := r.AttrStrings("srcs")
		if len(srcs) > 0 {
			newSrcs := make([]string, 0)
			for _, src := range srcs {
				if src == "package.proto" {
					continue
				}
				newSrcs = append(newSrcs, src)
			}
			r.SetAttr("srcs", protoc.DeduplicateAndSort(srcs))
			// log.Printf("added package.proto to %s //%s:%s", r.Kind(), args.Rel, r.Name())
			// deps := append(r.AttrStrings("deps"), "//thirdparty/protobuf/scalapb:scalapb_proto")
			// r.SetAttr("deps", protoc.DeduplicateAndSort(deps))
		}
	}

	rules := pkg.Rules()
	// empty := pkg.Empty()

	imports := make([]interface{}, len(rules))
	for i, r := range rules {
		imports[i] = r.PrivateAttr(config.GazelleImportsKey)
	}

	if debug && args.File != nil {
		log.Println("visited", args.Rel)
	}

	return language.GenerateResult{
		Gen: rules,
		// Empty:   empty,
		Imports: imports,
	}
}

// Imports returns a list of ImportSpecs that can be used to import the rule r.
// This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
func (sl *scalaLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	from := label.New("", f.Pkg, r.Name())

	pkg, ok := sl.packages[from.Pkg]
	if !ok {
		// log.Println("scala.Imports(): Unknown package", from.Pkg)
		return nil
	}

	provider := pkg.ruleProvider(r)
	// NOTE: gazelle attempts to index rules found in the build file regardless
	// of whether we returned the rule from GenerateRules or not, so this will
	// be nil in that case.
	if provider == nil {
		// log.Println("scala.Imports(): Unknown provider", from)
		return nil
	}

	return provider.Imports(c, r, f)
}

// Embeds returns a list of labels of rules that the given rule embeds. If a
// rule is embedded by another importable rule of the same language, only the
// embedding rule will be indexed. The embedding rule will inherit the imports
// of the embedded rule. Since SkyLark doesn't support embedding this should
// always return nil.
func (*scalaLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

func (sl *scalaLang) onResolvePhase() error {
	// log.Panicln("stopping at resolve phase")

	for _, r := range sl.resolvers {
		if gptl, ok := r.(GazellePhaseTransitionListener); ok {
			if err := gptl.OnResolvePhase(); err != nil {
				return err
			}
		}
	}
	if err := sl.scalaCompiler.OnResolvePhase(); err != nil {
		return err
	}

	// gather proto imports
	for from, imports := range protoc.GlobalResolver().Provided(ScalaLangName, ScalaLangName) {
		sl.importRegistry.Provides(from, imports)
	}

	// gather 1p/3p imports
	for _, rslv := range sl.resolvers {
		if ip, ok := rslv.(protoc.ImportProvider); ok {
			for from, imports := range ip.Provided(ScalaLangName, ScalaLangName) {
				sl.importRegistry.Provides(from, imports)
			}
		}
	}

	sl.importRegistry.OnResolve()
	return nil
}

// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each rule
// generated by language.GenerateRules in language.GenerateResult.Imports.
// Resolve generates a "deps" attribute (or the appropriate language-specific
// equivalent) for each import according to language-specific rules and
// heuristics.
func (sl *scalaLang) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	importsRaw interface{},
	from label.Label,
) {
	if !sl.isResolvePhase {
		sl.isResolvePhase = true
		if err := sl.onResolvePhase(); err != nil {
			log.Fatal(err)
		}
	}

	if pkg, ok := sl.packages[from.Pkg]; ok {
		provider := pkg.ruleProvider(r)
		if provider == nil {
			log.Printf("no known rule provider for %v", from)
		}
		if imports, ok := importsRaw.([]string); ok {
			provider.Resolve(c, ix, r, imports, from)
		} else {
			log.Printf("warning: resolve scala imports: expected []string, got %T", importsRaw)
		}
	} else {
		log.Printf("no known rule package for %v", from.Pkg)
	}
}

// CrossResolve calls all known resolvers and returns the first non-empty result.
func (sl *scalaLang) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	for _, r := range sl.resolvers {
		if result := r.CrossResolve(c, ix, imp, lang); len(result) > 0 {
			return result
		}
	}
	// final / fallback resolver.
	if result := sl.importRegistry.CrossResolve(c, ix, imp, lang); len(result) > 0 {
		return result
	}
	return nil
}

func hasPackageProto(files []string) bool {
	for _, f := range files {
		if filepath.Base(f) == "package.proto" {
			return true
		}
	}
	return false
}
