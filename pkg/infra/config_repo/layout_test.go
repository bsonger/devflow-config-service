package config_repo

import (
	"path/filepath"
	"testing"
)

func TestResolveLayoutFromServicePathAndEnv(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "applications/devflow-platform/services/devflow-app-service", "staging")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	if resolved.Env != "staging" {
		t.Fatalf("Env = %q, want %q", resolved.Env, "staging")
	}
	if got := filepath.ToSlash(resolved.Files[len(resolved.Files)-1]); got != "environments/staging.yaml" {
		t.Fatalf("last file = %q", got)
	}
}

func TestResolveLayoutFromEnvironmentPath(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "applications/devflow-platform/services/devflow-app-service/environments/prod", "")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	if resolved.Env != "prod" {
		t.Fatalf("Env = %q, want %q", resolved.Env, "prod")
	}
}
