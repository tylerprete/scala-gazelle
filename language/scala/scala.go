package scala

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/stackb/rules_proto/pkg/protoc"

	"github.com/stackb/scala-gazelle/pkg/progress"
)

const (
	ScalaLangName = "scala"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	var importRegistry *importRegistry
	depends := func(src, dst, kind string) {
		importRegistry.AddDependency(src, dst, kind)
	}
	packages := make(map[string]*scalaPackage)
	sourceResolver := newScalaSourceIndexResolver(depends)
	classResolver := newScalaClassIndexResolver(depends)
	mavenResolver := newMavenResolver()
	scalaCompiler := newScalaCompiler()
	// var scalaCompiler *scalaCompiler
	importRegistry = newImportRegistry(sourceResolver, classResolver, scalaCompiler)
	vizServer := newGraphvizServer(packages, importRegistry)

	out := progress.NewOut(os.Stderr)

	return &scalaLang{
		ruleRegistry:    globalRuleRegistry,
		scalaFileParser: sourceResolver,
		scalaCompiler:   scalaCompiler,
		packages:        packages,
		importRegistry:  importRegistry,
		resolvers: []ConfigurableCrossResolver{
			sourceResolver,
			classResolver,
			mavenResolver,
		},
		progress: progress.NewProgressOutput(out),
		viz:      vizServer,
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
	// has completed and deps resolution phase has started (it calls
	// onResolvePhase).
	isResolvePhase bool
	// resolvers is a list of cross resolver implementations.  Typically there
	// are two: one to help with third-party code, one to help with first-party
	// code.
	resolvers []ConfigurableCrossResolver
	// viz is the dependency vizualization engine
	viz *graphvizServer
	// lastPackage tracks if this is the last generated package
	lastPackage *scalaPackage
	// totalPackageCount is used for progress
	totalPackageCount int
	// remainingRules is a counter that tracks when all rules have been resolved.
	remainingRules int
	// totalRules is used for progress
	totalRules int
	// progress is the progress interface
	progress progress.Output
}

// Name implements part of the language.Language interface
func (sl *scalaLang) Name() string { return ScalaLangName }

// RegisterFlags implements part of the language.Language interface
func (sl *scalaLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	getOrCreateScalaConfig(c) // ignoring return value, only want side-effect

	for _, r := range sl.resolvers {
		r.RegisterFlags(fs, cmd, c)
	}

	sl.scalaCompiler.RegisterFlags(fs, cmd, c)
	sl.viz.RegisterFlags(fs, cmd, c)

	fs.IntVar(&sl.totalPackageCount, "total_package_count", 0, "number of total packages for the workspace (used for progress estimation)")
}

// CheckFlags implements part of the language.Language interface
func (sl *scalaLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	for _, r := range sl.resolvers {
		if err := r.CheckFlags(fs, c); err != nil {
			return err
		}
	}
	if err := sl.scalaCompiler.CheckFlags(fs, c); err != nil {
		return err
	}
	if err := sl.viz.CheckFlags(fs, c); err != nil {
		return err
	}
	return nil
}

// KnownDirectives implements part of the language.Language interface
func (*scalaLang) KnownDirectives() []string {
	return []string{
		ruleDirective,
		overrideDirective,
		indirectDependencyDirective,
		scalaExplainDependencies,
		mapKindImportNameDirective,
	}
}

// Loads implements part of the language.Language interface
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

// Configure implements part of the language.Language interface
func (sl *scalaLang) Configure(c *config.Config, rel string, f *rule.File) {
	if f == nil {
		return
	}
	if err := getOrCreateScalaConfig(c).parseDirectives(rel, f.Directives); err != nil {
		log.Fatalf("error while parsing rule directives in package %q: %v", rel, err)
	}
}

// Kinds implements part of the language.Language interface
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

// Fix implements part of the language.Language interface
func (sl *scalaLang) Fix(c *config.Config, f *rule.File) {
}

// GenerateRules implements part of the language.Language interface
func (sl *scalaLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {

	if sl.totalPackageCount > 0 {
		sl.progress.WriteProgress(progress.Progress{
			ID:      "walk",
			Action:  "generating rules",
			Total:   int64(sl.totalPackageCount),
			Current: int64(len(sl.packages)),
			Units:   "packages",
		})
	}

	cfg := getOrCreateScalaConfig(args.Config)

	pkg := newScalaPackage(sl.ruleRegistry, sl.scalaFileParser, sl.importRegistry, args.Rel, args.File, cfg)
	// search for child packages, but only assign if a parent has not already
	// been assigned.  Given that gazelle uses a DFS walk, we should assign the
	// child to the nearest parent.
	for rel, child := range sl.packages {
		if child.parent != nil {
			continue
		}
		if !strings.HasPrefix(rel, args.Rel) {
			continue
		}
		child.parent = pkg
		sl.importRegistry.AddDependency("pkg/"+args.Rel, "pkg/"+rel, "pkg")
	}
	sl.packages[args.Rel] = pkg
	sl.importRegistry.AddDependency("ws/default", "pkg/"+args.Rel, "ws")
	sl.lastPackage = pkg

	rules := pkg.Rules()
	sl.remainingRules += len(rules)
	// empty := pkg.Empty()

	imports := make([]interface{}, len(rules))
	for i, r := range rules {
		imports[i] = r.PrivateAttr(config.GazelleImportsKey)
		sl.importRegistry.AddDependency("pkg/"+args.Rel, "rule/"+label.New("", args.Rel, r.Name()).String(), "rule")
	}

	return language.GenerateResult{
		Gen: rules,
		// Empty:   empty,
		Imports: imports,
	}
}

// Imports implements part of the language.Language interface
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

// Embeds implements part of the language.Language interface
func (*scalaLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

// Resolve implements part of the language.Language interface
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
		sl.onResolve()
		sl.totalRules = sl.remainingRules
	}

	if pkg, ok := sl.packages[from.Pkg]; ok {
		pkg.Resolve(c, ix, rc, r, importsRaw, from)

		sl.remainingRules--

		if sl.remainingRules == 0 {
			sl.onEnd()
		}
	} else {
		log.Printf("no known rule package for %v", from.Pkg)
	}

	if sl.totalRules > 0 {
		sl.progress.WriteProgress(progress.Progress{
			ID:      "resolve",
			Action:  "resolving dependencies",
			Total:   int64(sl.totalRules),
			Current: int64(sl.totalRules - sl.remainingRules),
			Units:   "rules",
		})
	}
}

// CrossResolve implements part of the resolve.CrossResolver interface
func (sl *scalaLang) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	for _, r := range sl.resolvers {
		if result := r.CrossResolve(c, ix, imp, lang); len(result) > 0 {
			// log.Printf("scala.CrossResolve hit %T %s", r, imp.Imp)
			return result
		}
	}
	if result := sl.importRegistry.CrossResolve(c, ix, imp, lang); len(result) > 0 {
		// log.Printf("scala.CrossResolve hit %T %s", sl.importRegistry, imp.Imp)
		return result
	}
	return nil
}

// onResolve is called when gazelle transitions from the generate phase to the resolve phase
func (sl *scalaLang) onResolve() {

	for _, r := range sl.resolvers {
		if l, ok := r.(GazellePhaseTransitionListener); ok {
			l.OnResolve()
		}
	}

	sl.scalaCompiler.OnResolve()
	sl.viz.OnResolve()

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
}

// onEnd is called when the last rule has been resolved.
func (sl *scalaLang) onEnd() {
	sl.scalaCompiler.stop()
	// sl.recordDeps()
	sl.viz.OnEnd()
}

// recordDeps writes deps info to the graph once all rules resolved.
func (sl *scalaLang) recordDeps() {
	for _, pkg := range sl.packages {
		for _, r := range pkg.rules {
			from := label.New("", pkg.rel, r.Name())
			for _, dep := range r.AttrStrings("deps") {
				to, err := label.Parse(dep)
				if err != nil {
					continue
				}
				sl.importRegistry.AddDependency("rule/"+from.String(), "rule/"+to.String(), "depends")
			}
		}
	}
}

func hasPackageProto(files []string) bool {
	for _, f := range files {
		if filepath.Base(f) == "package.proto" {
			return true
		}
	}
	return false
}
