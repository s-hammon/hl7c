package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrMissingVersion = errors.New("no version number")
var ErrNoModels = errors.New("no models")

type versionSetting struct {
	Number string `yaml:"version"`
}

type Config struct {
	Meta   Meta         `yaml:"meta"`
	Models []Model      `yaml:"models"`
	Types  []CustomType `yaml:"types"`
}

func (c *Config) Compile(pkg string) string {
	var sb strings.Builder

	sb.WriteString("package " + pkg + "\n\n")

	if len(c.Meta.Imports) > 0 {
		sb.WriteString("import (\n")
		for _, imp := range c.Meta.Imports {
			sb.WriteString("\t\"" + imp + "\"\n")
		}
		sb.WriteString(")\n\n")
	}

	sb.WriteString(createBaseStruct())
	sb.WriteString(makeTypes(c.Types))
	sb.WriteString(makeModels(c.Models))

	return sb.String()
}

type Meta struct {
	Package string   `yaml:"package"`
	Imports []string `yaml:"imports"`
}

type Model struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

func (m *Model) createCustomUnmarshalJSON() string {
	sb := strings.Builder{}
	customFields := map[string]string{}

	// func signature
	sb.WriteString(fmt.Sprintf("func (m *%s) UnmarshalJSON(b []byte) error {\n", m.Name))
	// JSON alias
	sb.WriteString(fmt.Sprintf("\ttype alias %s\n", m.Name))
	sb.WriteString("\taux := &struct {\n")
	for _, field := range m.Fields {
		if mapType(field.Type) == "time.Time" {
			fmtTime := "20060102150405"
			if field.Type == "date" {
				fmtTime = "20060102"
			}
			customFields[field.Name] = fmt.Sprintf("time.Parse(\"%s\", aux.%s)", fmtTime, field.Name)
			sb.WriteString(fmt.Sprintf("\t\t%s string `json:\"%s\"`\n", field.Name, field.Tag))
		}
	}
	sb.WriteString("\t\t*alias\n\t}{\n")
	sb.WriteString("\t\talias: (*alias)(m),\n\t}\n")

	// Unmarshal
	sb.WriteString("\tif err := json.Unmarshal(b, &aux); err != nil {\n\t\treturn err \n\t}\n")
	// Custom fields (rework, won't always be time.Time)
	for k, v := range customFields {
		sb.WriteString(fmt.Sprintf("\tm.%s, _ = %s\n", k, v))
	}
	sb.WriteString(fmt.Sprintln())

	// Base fields
	sb.WriteString("\tm.ID = uuid.New()\n\tm.CreatedAt = time.Now()\n\tm.UpdatedAt = time.Now()\n")
	sb.WriteString("\treturn nil\n}\n\n")

	return sb.String()
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Tag  string `yaml:"tag"`
}

func (f *Field) writeField() string {
	s := fmt.Sprintf(
		"\t%s %s %s\n",
		f.Name,
		mapType(f.Type),
		toGoTag(f.Tag),
	)

	return s
}

type CustomType struct {
	Name                string  `yaml:"name"`
	Fields              []Field `yaml:"fields"`
	needCustomUnmarshal bool
}

func (c *CustomType) createCustomUnmarshalJSON() string {
	sb := strings.Builder{}
	customFields := map[string]string{}

	sb.WriteString(fmt.Sprintf("func (t *%s) UnmarshalJSON(b []byte) error {\n", c.Name))
	sb.WriteString(fmt.Sprintf("\ttype alias %s\n", c.Name))
	sb.WriteString("\taux := &struct {\n")
	for _, field := range c.Fields {
		if mapType(field.Type) == "time.Time" {
			fmtTime := "20060102150405"
			if field.Type == "date" {
				fmtTime = "20060102"
			}
			customFields[field.Name] = fmt.Sprintf("time.Parse(\"%s\", aux.%s)", fmtTime, field.Name)
			sb.WriteString(fmt.Sprintf("\t\t%s string `json:\"%s\"`\n", field.Name, field.Tag))
		}
	}
	sb.WriteString("\t\t*alias\n\t}{\n")
	sb.WriteString("\t\talias: (*alias)(t),\n\t}\n")

	sb.WriteString("\tif err := json.Unmarshal(b, &aux); err != nil {\n\t\treturn err \n\t}\n")
	for k, v := range customFields {
		sb.WriteString(fmt.Sprintf("\tt.%s, _ = %s\n", k, v))
	}

	sb.WriteString("\treturn nil\n}\n\n")

	return sb.String()
}

func makeTypes(types []CustomType) string {
	var sb strings.Builder
	for _, t := range types {
		sb.WriteString(fmt.Sprintf("type %s struct {\n", t.Name))

		for _, field := range t.Fields {
			sb.WriteString(field.writeField())

		}
		sb.WriteString("}\n\n")

		if t.needCustomUnmarshal {
			sb.WriteString(t.createCustomUnmarshalJSON())
		}
	}

	return sb.String()
}

func makeModels(models []Model) string {
	var sb strings.Builder
	for _, m := range models {
		sb.WriteString(fmt.Sprintf("type %s struct {\n", m.Name))
		sb.WriteString("\tBase\n")

		for _, field := range m.Fields {
			sb.WriteString(field.writeField())
		}
		sb.WriteString("}\n\n")

		sb.WriteString(m.createCustomUnmarshalJSON())
	}

	return sb.String()
}

func ParseConfig(rd io.Reader) (Config, error) {
	var buf bytes.Buffer
	var config Config
	var version versionSetting

	ver := io.TeeReader(rd, &buf)
	dec := yaml.NewDecoder(ver)
	if err := dec.Decode(&version); err != nil {
		return config, err
	}

	return parseConfig(&buf)
}

func parseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)
	dec.KnownFields(true)

	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return conf, err
	}
	if len(conf.Models) == 0 {
		return conf, ErrNoModels
	}

	for _, module := range []string{"encoding/json", "github.com/google/uuid"} {
		if !contains(conf.Meta.Imports, module) {
			conf.Meta.Imports = append(conf.Meta.Imports, module)
		}
	}

	if err := conf.Sanitize(); err != nil {
		return Config{}, err
	}

	return conf, nil
}
