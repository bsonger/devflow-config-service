package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestConfigurationContract(t *testing.T) {
	typ := reflect.TypeOf(Configuration{})
	for _, field := range []string{"ApplicationID", "Name", "Env", "LatestRevisionNo", "LatestRevisionID"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("Configuration missing field %s", field)
		}
		if field == "ApplicationID" && f.Type != reflect.TypeOf(uuid.UUID{}) {
			t.Fatalf("Configuration.ApplicationID type = %v, want uuid.UUID", f.Type)
		}
	}
	if _, ok := typ.FieldByName("Files"); ok {
		t.Fatal("Configuration should not own mutable files directly")
	}
}

func TestConfigurationRevisionContract(t *testing.T) {
	typ := reflect.TypeOf(ConfigurationRevision{})
	for _, field := range []string{"ConfigurationID", "RevisionNo", "Files", "ContentHash", "CreatedAt"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("ConfigurationRevision missing field %s", field)
		}
		if field == "ConfigurationID" && f.Type != reflect.TypeOf(uuid.UUID{}) {
			t.Fatalf("ConfigurationRevision.ConfigurationID type = %v, want uuid.UUID", f.Type)
		}
	}
}

func TestBaseModelWithCreateDefault(t *testing.T) {
	var base BaseModel
	base.WithCreateDefault()

	if base.ID == uuid.Nil {
		t.Fatal("BaseModel.WithCreateDefault should assign a UUID")
	}
	if base.CreatedAt.IsZero() || base.UpdatedAt.IsZero() {
		t.Fatal("BaseModel.WithCreateDefault should set timestamps")
	}
}
