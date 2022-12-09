package scala

import (
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	sppb "github.com/stackb/scala-gazelle/build/stack/gazelle/scala/parse"
	"github.com/stackb/scala-gazelle/pkg/collections"
	"github.com/stackb/scala-gazelle/pkg/resolver"
)

type scalaRule struct {
	// Rule is an embedded struct (FIXME: why embed this?).
	*rule.Rule
	// the parent config
	scalaConfig *scalaConfig
	// from is the label for the rule.
	from label.Label
	// files that are included in the rule.
	files []*sppb.File
	// the import resolver to which we chain to when self-imports are not matched.
	next resolver.SymbolResolver
	// the ruleScope implementation to which we provide symbols.
	ruleScope resolver.Scope
	// the scope from which we lookup both local and global imports.
	importScope resolver.Scope
	// extendedTypes is a mapping from the required type to the symbol that
	// needs it. for example, if 'class Foo extendedTypes Bar', "Bar" is the map
	// key and "Foo" will be the value.
	extendedTypes map[string][]string
	// exports represent symbols that are importable by other rules.
	exports map[string]resolve.ImportSpec
	// scope is a map of symbols that are in scope.  For the import
	// 'com.foo.Bar', the map key is 'Bar' and the map value is the symbol for
	// it.
	scope resolver.SymbolMap
}

func newScalaRule(
	scalaConfig *scalaConfig,
	next resolver.SymbolResolver,
	parentScope resolver.Scope,
	localScope resolver.Scope,
	r *rule.Rule,
	from label.Label,
	files []*sppb.File,
) *scalaRule {
	scalaRule := &scalaRule{
		Rule:          r,
		scalaConfig:   scalaConfig,
		from:          from,
		files:         files,
		next:          next,
		ruleScope:     localScope,
		importScope:   resolver.NewChainScope(localScope, parentScope),
		extendedTypes: make(map[string][]string),
		exports:       make(map[string]resolve.ImportSpec),
	}
	scalaRule.visitFiles()
	return scalaRule
}

// ResolveSymbol implements the resolver.SymbolResolver interface
func (r *scalaRule) ResolveSymbol(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, imp string) (*resolver.Symbol, error) {
	if symbol, ok := r.ruleScope.GetSymbol(imp); ok {
		return symbol, nil
	}
	return r.next.ResolveSymbol(c, ix, from, lang, imp)
}

// Imports implements part of the scalarule.Rule interface.
func (r *scalaRule) Imports() resolver.ImportMap {
	imports := resolver.NewImportMap()
	impLang := scalaLangName

	// direct
	for _, file := range r.files {
		for _, imp := range file.Imports {
			imports.Put(resolver.NewDirectImport(imp, file))
		}
	}

	// if this rule has a main_class
	if mainClass := r.AttrString("main_class"); mainClass != "" {
		imports.Put(resolver.NewMainClassImport(mainClass))
	}

	// add import required from extends clauses
	scope := r.getOrCreateImportScope()
	for imp, src := range r.extendedTypes {
		// check if the import is actually a symbol in scope.  If yes, use the
		// fully-qualified import name.
		if symbol, ok := scope.Get(imp); ok {
			imp = symbol.Name
		}
		imports.Put(resolver.NewExtendsImport(imp, src[0])) // use first occurrence as source arg
	}

	// Initialize a list of symbols to find implicits for from all known
	// imports. Include all symbols that are defined in the rule too (a
	// gazelle:resolve_with directive should apply to them too).
	required := collections.StringStack(imports.Keys())
	for _, export := range r.Exports() {
		required = append(required, export.Imp)
	}

	// Gather implicit imports transitively.
	for !required.IsEmpty() {
		src, _ := required.Pop()
		for _, dst := range r.scalaConfig.getImplicitImports(impLang, src) {
			required.Push(dst)
			imports.Put(resolver.NewImplicitImport(dst, src))
		}
	}

	return imports
}

// Exports implements part of the scalarule.Rule interface.
func (r *scalaRule) Exports() []resolve.ImportSpec {
	exports := make([]resolve.ImportSpec, 0, len(r.exports))
	for _, v := range r.exports {
		exports = append(exports, v)
	}

	sort.Slice(exports, func(i, j int) bool {
		a := exports[i]
		b := exports[j]
		return a.Imp < b.Imp
	})

	return exports
}

// Files implements part of the scalarule.Rule interface.
func (r *scalaRule) Files() []*sppb.File {
	return r.files
}

func (r *scalaRule) visitFiles() {
	for _, file := range r.files {
		r.visitFile(file)
	}
}

func (r *scalaRule) visitFile(file *sppb.File) {
	for _, imp := range file.Classes {
		r.putSymbol(imp, sppb.ImportType_CLASS)
		r.putExport(imp)
	}
	for _, imp := range file.Objects {
		r.putSymbol(imp, sppb.ImportType_OBJECT)
		r.putExport(imp)
	}
	for _, imp := range file.Traits {
		r.putSymbol(imp, sppb.ImportType_TRAIT)
		r.putExport(imp)
	}
	for _, imp := range file.Types {
		r.putSymbol(imp, sppb.ImportType_TYPE)
		r.putExport(imp)
	}
	for _, imp := range file.Vals {
		r.putSymbol(imp, sppb.ImportType_VALUE)
		r.putExport(imp)
	}
	for token, extends := range file.Extends {
		r.visitExtends(token, extends)
	}
}

func (r *scalaRule) visitExtends(token string, types *sppb.ClassList) {
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 {
		log.Fatalf("invalid extends token: %q: should have form '(class|interface|object) com.foo.Bar' ", token)
	}

	kind := parts[0]
	symbol := parts[1]

	r.visitKindExtends(kind, symbol, types)
}

func (r *scalaRule) visitKindExtends(kind, symbol string, types *sppb.ClassList) {
	switch kind {
	case "class":
		r.visitClassExtends(symbol, types)
	case "interface":
		r.visitInterfaceExtends(symbol, types)
	case "object":
		r.visitObjectExtends(symbol, types)
	}
}

func (r *scalaRule) visitClassExtends(imp string, types *sppb.ClassList) {
	r.putExtendedTypes(imp, types)
}

func (r *scalaRule) visitInterfaceExtends(imp string, types *sppb.ClassList) {
	r.putExtendedTypes(imp, types)
}

func (r *scalaRule) visitObjectExtends(imp string, types *sppb.ClassList) {
	r.putExtendedTypes(imp, types)
}

func (r *scalaRule) putExtendedTypes(imp string, types *sppb.ClassList) {
	for _, dst := range types.Classes {
		r.putExtendedType(imp, dst)
	}
}

func (r *scalaRule) putExtendedType(src, dst string) {
	r.extendedTypes[dst] = append(r.extendedTypes[dst], src)
}

func (r *scalaRule) putExport(imp string) {
	r.exports[imp] = resolve.ImportSpec{Imp: imp, Lang: scalaLangName}
}

func (r *scalaRule) putSymbol(imp string, impType sppb.ImportType) {
	// since we don't need to resolve same-rule symbols to a different label,
	// record all imports as label.NoLabel!
	r.ruleScope.PutSymbol(resolver.NewSymbol(impType, imp, "self-import", label.NoLabel))
}

func (r *scalaRule) getOrCreateImportScope() resolver.SymbolMap {
	if r.scope == nil {
		r.scope = make(resolver.SymbolMap)
		r.visitImports()
	}
	return r.scope
}

func (r *scalaRule) visitImports() {
	for _, file := range r.files {
		for _, imp := range file.Imports {
			r.visitImport(imp)
		}
	}
}

func (r *scalaRule) visitImport(imp string) {
	if prefix, ok := isWildcardImport(imp); ok {
		r.visitWildcardImport(prefix)
	} else {
		r.visitExplicitImport(imp)
	}
}

func (r *scalaRule) visitExplicitImport(imp string) {
	if symbol, ok := r.importScope.GetSymbol(imp); ok {
		r.putSymbolInScope(symbol)
	}
}

func (r *scalaRule) visitWildcardImport(prefix string) {
	for _, symbol := range r.importScope.GetSymbols(prefix) {
		r.putSymbolInScope(symbol)
	}
}

func (r *scalaRule) putSymbolInScope(symbol *resolver.Symbol) {
	r.scope.Add(symbol)
}

func isWildcardImport(imp string) (string, bool) {
	if !strings.HasSuffix(imp, "._") {
		return "", false
	}
	return imp[:len(imp)-len("._")], true
}
