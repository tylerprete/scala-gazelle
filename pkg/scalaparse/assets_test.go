package scalaparse

import (
	"testing"
)

func TestEmbed(t *testing.T) {
	if len(scalaparserMjs) == 0 {
		t.Error("embedded scalaparser.mjs script is missing")
	}
	if len(scalametaParsersIndexJs) == 0 {
		t.Error("embedded node_modules/scalameta-parsers/index.js script is missing")
	}
	if len(nodeExe) == 0 {
		t.Errorf("embedded node.exe is missing: len %d", len(nodeExe))
	}
}
