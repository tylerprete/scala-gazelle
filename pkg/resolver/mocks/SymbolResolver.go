// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	config "github.com/bazelbuild/bazel-gazelle/config"
	label "github.com/bazelbuild/bazel-gazelle/label"

	mock "github.com/stretchr/testify/mock"

	resolve "github.com/bazelbuild/bazel-gazelle/resolve"

	resolver "github.com/stackb/scala-gazelle/pkg/resolver"
)

// SymbolResolver is an autogenerated mock type for the SymbolResolver type
type SymbolResolver struct {
	mock.Mock
}

// ResolveSymbol provides a mock function with given fields: c, ix, from, lang, sym
func (_m *SymbolResolver) ResolveSymbol(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, sym string) (*resolver.Symbol, error) {
	ret := _m.Called(c, ix, from, lang, sym)

	var r0 *resolver.Symbol
	if rf, ok := ret.Get(0).(func(*config.Config, *resolve.RuleIndex, label.Label, string, string) *resolver.Symbol); ok {
		r0 = rf(c, ix, from, lang, sym)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resolver.Symbol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*config.Config, *resolve.RuleIndex, label.Label, string, string) error); ok {
		r1 = rf(c, ix, from, lang, sym)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSymbolResolver interface {
	mock.TestingT
	Cleanup(func())
}

// NewSymbolResolver creates a new instance of SymbolResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSymbolResolver(t mockConstructorTestingTNewSymbolResolver) *SymbolResolver {
	mock := &SymbolResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
