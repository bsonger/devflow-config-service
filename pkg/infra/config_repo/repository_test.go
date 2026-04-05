package config_repo

import (
	"context"
	"path/filepath"
	"testing"
)

func TestRepositoryReadSnapshot(t *testing.T) {
	repo := NewRepository(Options{
		RootDir:    filepath.Join("testdata", "config-repo"),
		DefaultRef: "main",
	})

	snapshot, err := repo.ReadSnapshot(context.Background(), "applications/example-app/staging")
	if err != nil {
		t.Fatalf("ReadSnapshot returned error: %v", err)
	}
	if snapshot.SourceCommit != "main" {
		t.Fatalf("SourceCommit = %q, want %q", snapshot.SourceCommit, "main")
	}
	if len(snapshot.Files) != 2 {
		t.Fatalf("len(Files) = %d, want 2", len(snapshot.Files))
	}
	if snapshot.Files[0].Name != "app.yaml" {
		t.Fatalf("first file = %q, want %q", snapshot.Files[0].Name, "app.yaml")
	}
	if snapshot.SourceDigest == "" {
		t.Fatal("SourceDigest should not be empty")
	}
}
