package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Config
	}{
		{
			name: "test valid model",
			data: []byte(`
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
        type: timestamp
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
`),
			want: Config{
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
							{
								Name: "Mrn",
								Type: "CX",
								Tag:  "PID.3",
							},
							{
								Name: "Name",
								Type: "string",
								Tag:  "PID.5",
							},
							{
								Name: "Dob",
								Type: "timestamp",
								Tag:  "PID.7",
							},
						},
					},
				},
				Types: []CustomType{
					{
						Name: "CX",
						Fields: []Field{
							{"Id", "", "1"},
							{"CheckDigit", "", "2"},
							{"CheckDigitScheme", "", "3"},
							{"AssigningAuthority", "", "4"},
							{"IdentifierTypeCode", "", "5"},
							{"AssigningFacility", "", "6"},
						},
					},
				},
			},
		},
		{
			name: "test import missing modules",
			data: []byte(`
meta:
  package: objects

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
        type: timestamp
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
`),
			want: Config{
				Meta: Meta{
					Package: "objects",
					Imports: []string{
						"encoding/json",
						"github.com/google/uuid",
						"time",
					},
				},
				Models: []Model{
					{
						Name: "Patient",
						Fields: []Field{
							{
								Name: "Mrn",
								Type: "CX",
								Tag:  "PID.3",
							},
							{
								Name: "Name",
								Type: "string",
								Tag:  "PID.5",
							},
							{
								Name: "Dob",
								Type: "timestamp",
								Tag:  "PID.7",
							},
						},
					},
				},
				Types: []CustomType{
					{
						Name: "CX",
						Fields: []Field{
							{"Id", "", "1"},
							{"CheckDigit", "", "2"},
							{"CheckDigitScheme", "", "3"},
							{"AssigningAuthority", "", "4"},
							{"IdentifierTypeCode", "", "5"},
							{"AssigningFacility", "", "6"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want: %v, got %v", tt.want, got)
			}
		})
	}
}

func TestNewConfigError(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "invalid YAML",
			data: []byte(`invalid`),
		},
		{
			name: "missing models",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				    - time
			`),
		},
		{
			name: "type with no fields",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				    - time
				types:
				  - name: CX
			`),
		},
		{
			name: "type with no name",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				    - time
				types:
				  - fields:
					- name: id
					- name: checkDigit
			`),
		},
		{
			name: "model with no fields",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				    - time
				models:
				  - name: patient
			`),
		},
		{
			name: "model with no name",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				    - time
				models:
				  - fields:
				  	- name: name
					  type: string
					  tag: PID.5
			`),
		},
		{
			name: "model with undefined type",
			data: []byte(`
				meta:
				  package: objects
				  imports:
				  	- time
				models:
				  - name: patient
				    fields:
					  - name: mrn
					    type: CX
						tag: PID.3
			`),
		},
	}

	for _, tt := range tests {
		if _, err := New(tt.data); err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
	}
}

func TestUnmarshalCustomType(t *testing.T) {
	data := []byte(`
meta:
  package: objects

models:
  - name: patient
    fields:
      - name: account
        type: Account
        tag: PID.3

types:
  - name: Account
    fields:
      - name: id
      - name: createDate
        type: timestamp
`)

	config, err := New(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !config.Types[0].needCustomUnmarshal {
		t.Error("expected needCustomUnmarshal to be true")
	}
}
