package config

import (
	"bytes"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	validYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7

types:
  - name: CX
    fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
`

	typeNoFieldYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7

types:
  - name: CX
    fields:
`

	typeMissingNameYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7

types:
  - fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
`

	modelMissingFieldsYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient

types:
  - name: CX
    fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
`

	modelMissingNameYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7

types:
  - name: CX
    fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
`

	fieldMissingNameYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7

types:
  - name: CX
    fields:
      - name: id
      - name: checkDigit
      - name: checkDigitScheme
      - name: assigningAuthority
      - name: identifierTypeCode
      - name: assigningFacility
`

	undefinedTypeYAML = `
meta:
  package: objects
  imports:
    - time
    - "github.com/google/uuid"

models:
  - name: patient
    fields:
      - name: mrn
        type: CX
        tag: PID.3
      - name: name
        type: string
        tag: PID.5
      - name: dob
        type: date
        tag: PID.7
`
)

func TestSanitize(t *testing.T) {
	dec := yaml.NewDecoder(bytes.NewBufferString(validYAML))
	dec.KnownFields(true)

	var conf Config
	if err := dec.Decode(&conf); err != nil {
		t.Fatal(err)
	}

	if err := conf.Sanitize(); err != nil {
		t.Fatal(err)
	}

	if len(conf.Meta.Imports) != 2 {
		t.Errorf("expected 2 imports, got %d", len(conf.Meta.Imports))
	}
	if len(conf.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(conf.Models))
	}
	if len(conf.Types) != 1 {
		t.Errorf("expected 1 type, got %d", len(conf.Types))
	}

	expectedFields := []Field{
		{Name: "Mrn", Type: "CX", Tag: "PID.3"},
		{Name: "Name", Type: "string", Tag: "PID.5"},
		{Name: "Dob", Type: "date", Tag: "PID.7"},
	}

	if !reflect.DeepEqual(conf.Models[0].Fields, expectedFields) {
		t.Errorf("expected fields %+v, got %+v", expectedFields, conf.Models[0].Fields)
	}
}

func TestSanitizeError(t *testing.T) {
	tests := []struct {
		name   string
		yaml   string
		errMsg error
		errObj string
	}{
		{"type w/ no fields", typeNoFieldYAML, ErrTypeHasNoFields, "CX"},
		{"type missing name", typeMissingNameYAML, ErrTypeMissingName, ""},
		{"model missing fields", modelMissingFieldsYAML, ErrModelHasNoFields, "patient"},
		{"model missing name", modelMissingNameYAML, ErrModelMissingName, ""},
		{"field missing name", fieldMissingNameYAML, ErrFieldMissingName, "patient"},
		{"undefined type", undefinedTypeYAML, ErrUndefinedType, "CX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := yaml.NewDecoder(bytes.NewBufferString(tt.yaml))
			dec.KnownFields(true)

			var conf Config
			if err := dec.Decode(&conf); err != nil {
				t.Fatal(err)
			}

			err := conf.Sanitize()
			if err == nil {
				t.Fatalf("%s: expected error, got nil", tt.name)
			}

			wantErr := &CustomConfigError{
				Message: tt.errMsg,
				Obj:     tt.errObj,
			}

			if wantErr.Error() != err.Error() {
				t.Errorf("%s: expected error message to be %q, got %q", tt.name, wantErr.Error(), err.Error())
				return
			}
		})
	}
}

func TestMapType(t *testing.T) {
	yamlTypes := []struct {
		in   string
		want string
	}{
		{"", "string"},
		{"string", "string"},
		{"int", "int"},
		{"date", "time.Time"},
		{"timestamp", "time.Time"},
		{"uuid", "uuid.UUID"},
	}

	for _, yamlType := range yamlTypes {
		if got := mapType(yamlType.in); got != yamlType.want {
			t.Errorf("%q: expected %q, got %q", yamlType.in, yamlType.want, got)
		}
	}
}
