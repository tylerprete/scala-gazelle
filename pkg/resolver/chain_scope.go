package resolver

import (
	"fmt"
)

// ChainScope implements Scope over a chain of registries.
type ChainScope struct {
	chain []Scope
}

func NewChainScope(chain ...Scope) *ChainScope {
	return &ChainScope{
		chain: chain,
	}
}

// PutSymbol is not supported and will panic.
func (r *ChainScope) PutSymbol(known *Symbol) error {
	return fmt.Errorf("unsupported operation: PutSymbol")
}

// GetSymbol implements part of the Scope interface
func (r *ChainScope) GetSymbol(imp string) (*Symbol, bool) {
	for _, next := range r.chain {
		if known, ok := next.GetSymbol(imp); ok {
			return known, true
		}
	}
	return nil, false
}

// GetSymbols implements part of the Scope interface
func (r *ChainScope) GetSymbols(prefix string) []*Symbol {
	for _, next := range r.chain {
		if known := next.GetSymbols(prefix); len(known) > 0 {
			return known
		}
	}
	return nil
}
