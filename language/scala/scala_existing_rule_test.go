package scala

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/scala-gazelle/pkg/index"
)

// This looks important: https://github.com/sbt/zinc/blob/7c796ce65217096ce71be986149b2e769f8b33af/internal/zinc-core/src/main/scala/sbt/internal/inc/Relations.scala

func TestScalaExportSymbols(t *testing.T) {
	for name, tc := range map[string]struct {
		resolved index.ScalaFileSpec
		file     index.ScalaFileSpec
		want     []string
		wantErr  error
	}{
		"degenerate": {},
		"miss": {
			resolved: index.ScalaFileSpec{},
			file: index.ScalaFileSpec{
				Filename: "foo.scala",
				Extends: map[string][]string{
					"class trumid.common.akka.grpc.AbstractGrpcService": {
						"LazyLogging",
						"ReadinessReporter",
					},
				},
			},
			wantErr: fmt.Errorf(`failed to resolve name "LazyLogging" in file foo.scala!`),
		},
		"hit": {
			resolved: index.ScalaFileSpec{
				// contrived these would live in the same file
				Objects: []string{"com.typesafe.scalalogging.LazyLogging"},
				Traits:  []string{"com.foo.ReadinessReporter"},
			},
			file: index.ScalaFileSpec{
				Extends: map[string][]string{
					"class trumid.common.akka.grpc.AbstractGrpcService": {
						"LazyLogging",
						"ReadinessReporter",
					},
				},
			},
			want: []string{"com.typesafe.scalalogging.LazyLogging", "com.foo.ReadinessReporter"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolvers := []NameResolver{resolveNameInFile(&tc.resolved)}
			got, err := scalaExportSymbols(&tc.file, resolvers)
			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if diff := cmp.Diff(tc.wantErr.Error(), err.Error(), cmpopts.EquateErrors()); diff != "" {
					t.Errorf("error (-want +got):\n%s", diff)
				}
			} else {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestResolveNameInFile(t *testing.T) {
	for name, tc := range map[string]struct {
		file index.ScalaFileSpec
		name string
		want string
	}{
		"degenerate": {
			want: ``,
		},
		"miss": {
			file: index.ScalaFileSpec{},
			name: "Bar",
			want: "",
		},
		"hit trait": {
			file: index.ScalaFileSpec{Traits: []string{"com.foo.Bar"}},
			name: "Bar",
			want: "com.foo.Bar",
		},
		"hit class": {
			file: index.ScalaFileSpec{Classes: []string{"com.foo.Bar"}},
			name: "Bar",
			want: "com.foo.Bar",
		},
		"hit object": {
			file: index.ScalaFileSpec{Objects: []string{"com.foo.Bar"}},
			name: "Bar",
			want: "com.foo.Bar",
		},
		"hit type": {
			file: index.ScalaFileSpec{Types: []string{"com.foo.Bar"}},
			name: "Bar",
			want: "com.foo.Bar",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := resolveNameInFile(&tc.file)(tc.name)
			if tc.want == "" && !ok {
				return
			}
			if tc.want == "" && ok {
				t.Fatal("resolveNameInFile failed: expected miss")
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestResolveNameInFile (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveNameInLabelImportMap(t *testing.T) {
	for name, tc := range map[string]struct {
		resolved map[string]string
		name     string
		want     string
	}{
		"degenerate": {
			want: ``,
		},
		"miss": {
			name: "LazyLogging",
			want: "",
		},
		"hit": {
			resolved: map[string]string{
				"com.typesafe.scalalogging.LazyLogging": "@maven//:com_typesafe_scala_logging_scala_logging_2_12",
			},
			name: "LazyLogging",
			want: "com.typesafe.scalalogging.LazyLogging",
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolved := make(labelImportMap)
			for imp, origin := range tc.resolved {
				l, _ := label.Parse(origin)
				resolved.Set(l, imp)
			}
			got, ok := resolveNameInLabelImportMap(resolved)(tc.name)
			if tc.want == "" && !ok {
				return
			}
			if tc.want == "" && ok {
				t.Fatal("resolvedInLabelImportMap failed: expected miss")
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestResolvedInLabelImportMap (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMakeLabeledListExpr(t *testing.T) {
	for name, tc := range map[string]struct {
		// prelude is an optional chunk of BUILD file content
		directives []rule.Directive
		// resolved is a mapping from import -> label
		resolved map[string]string
		// want is the expected rule appearance
		want string
	}{
		"degenerate": {
			want: `testkind(
    name = "testname",
    deps = [],
)
`,
		},
		"simple": {
			resolved: map[string]string{
				"com.typesafe.scalalogging.LazyLogging": "@maven//:com_typesafe_scala_logging_scala_logging_2_12",
			},
			want: `testkind(
    name = "testname",
    deps = ["@maven//:com_typesafe_scala_logging_scala_logging_2_12"],
)
`,
		},
		"simple+reason": {
			directives: []rule.Directive{{"scala_explain_dependencies", "true"}},
			resolved: map[string]string{
				"com.typesafe.scalalogging.LazyLogging": "@maven//:com_typesafe_scala_logging_scala_logging_2_12",
			},
			want: `testkind(
    name = "testname",
    deps = [
        # com.typesafe.scalalogging.LazyLogging
        "@maven//:com_typesafe_scala_logging_scala_logging_2_12",
    ],
)
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := config.New()
			sc := getOrCreateScalaConfig(c)
			sc.parseDirectives("", tc.directives)
			from := label.New("", "pkg", "rule")
			resolved := make(labelImportMap)
			for imp, origin := range tc.resolved {
				l, _ := label.Parse(origin)
				resolved.Set(l, imp)
			}
			expr := makeLabeledListExpr(c, from, resolved)
			r := rule.NewRule("testkind", "testname")
			r.SetAttr("deps", expr)
			want := tc.want
			got := printRule(r)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestMakeLabeledListExpr (-want +got):\n%s", diff)
			}
		})
	}
}

func printRule(rules ...*rule.Rule) string {
	file := rule.EmptyFile("", "")
	for _, r := range rules {
		r.Insert(file)
	}
	return string(file.Format())
}
