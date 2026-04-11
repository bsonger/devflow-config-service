package config_repo

import (
	"path/filepath"
	"testing"
)

func TestResolveLayoutFromServicePathAndEnv(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "apps/11111111-1111-1111-1111-111111111111/staging/configmap", "staging")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	if got := filepath.ToSlash(resolved.SourcePath); got != "apps/11111111-1111-1111-1111-111111111111/staging/configmap" {
		t.Fatalf("SourcePath = %q", got)
	}
	if got := filepath.ToSlash(resolved.Files[len(resolved.Files)-1]); got != "logging.yaml" {
		t.Fatalf("last file = %q", got)
	}
}

func TestResolveLayoutFromEnvironmentPath(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "apps/11111111-1111-1111-1111-111111111111/staging/configmap", "")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	if len(resolved.Files) != 2 {
		t.Fatalf("len(Files) = %d", len(resolved.Files))
	}
}
