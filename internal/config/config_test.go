package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	r := strings.NewReader(validYAML)
	got, err := ParseConfig(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Meta.Package != "objects" {
		t.Errorf("expected 'objects', got '%s'", got.Meta.Package)
	}

	if len(got.Meta.Imports) != 3 {
		t.Errorf("expected 3 imports, got %d", len(got.Meta.Imports))
	}
	if len(got.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(got.Models))
	}

	want := Config{
		Meta: Meta{
			Package: "objects",
			Imports: []string{
				"time",
				"github.com/google/uuid",
				"encoding/json",
			},
		},
		Models: []Model{
			{
				Name: "Patient",
				Fields: []Field{
					{Name: "Mrn", Type: "CX", Tag: "PID.3"},
					{Name: "Name", Type: "string", Tag: "PID.5"},
					{Name: "Dob", Type: "date", Tag: "PID.7"},
				},
			},
		},
		Types: []CustomType{
			{
				Name: "CX",
				Fields: []Field{
					{Name: "Id", Type: "", Tag: "1"},
					{Name: "CheckDigit", Type: "", Tag: "2"},
					{Name: "CheckDigitScheme", Type: "", Tag: "3"},
					{Name: "AssigningAuthority", Type: "", Tag: "4"},
					{Name: "IdentifierTypeCode", Type: "", Tag: "5"},
					{Name: "AssigningFacility", Type: "", Tag: "6"},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %+v, got %+v", want, got)
	}
}

func TestParseConfigError(t *testing.T) {
	badYAML := `invalid_yaml`
	r := strings.NewReader(badYAML)

	_, err := parseConfig(r)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestFieldWriteField(t *testing.T) {
	tests := []struct {
		field Field
		want  string
	}{
		{
			Field{Name: "Mrn", Type: "CX", Tag: "PID.3"},
			"\tMrn CX `json:\"PID.3\"`\n",
		},
		{
			Field{Name: "Name", Type: "string"},
			"\tName string `json:\"\"`\n",
		},
	}

	for _, tt := range tests {
		got := tt.field.writeField()
		if got != tt.want {
			t.Errorf("expected '%s', got '%s'", tt.want, got)
		}
	}
}
