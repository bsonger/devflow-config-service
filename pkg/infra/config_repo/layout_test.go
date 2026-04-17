package config_repo

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveLayoutFromServicePathAndEnv(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "applications/devflow-platform/services/devflow-app-service", "staging")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	if got := filepath.ToSlash(resolved.SourcePath); got != "applications/devflow-platform/services/devflow-app-service" {
		t.Fatalf("SourcePath = %q", got)
	}
	want := []string{
		"configuration.yaml",
		"environments/staging.yaml",
	}
	if !reflect.DeepEqual(resolved.Files, want) {
		t.Fatalf("Files = %#v, want %#v", resolved.Files, want)
	}
}

func TestResolveLayoutFromServicePathBaseFilesOnly(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "applications/devflow-platform/services/devflow-app-service", "base")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	want := []string{
		"configuration.yaml",
	}
	if !reflect.DeepEqual(resolved.Files, want) {
		t.Fatalf("Files = %#v, want %#v", resolved.Files, want)
	}
}

func TestResolveLayoutLegacyFlatDirectoryStillWorks(t *testing.T) {
	resolved, err := resolveLayout(filepath.Join("testdata", "config-repo"), "apps/11111111-1111-1111-1111-111111111111/staging/configmap", "")
	if err != nil {
		t.Fatalf("resolveLayout returned error: %v", err)
	}
	want := []string{"app.yaml", "logging.yaml"}
	if !reflect.DeepEqual(resolved.Files, want) {
		t.Fatalf("Files = %#v, want %#v", resolved.Files, want)
	}
}
