package config

import (
	"errors"
	"log"
	"strconv"
)

type CustomConfigError struct {
	Message error
	Obj     string
}

func (e *CustomConfigError) Error() string {
	if e.Obj == "" {
		return e.Message.Error()
	}
	return e.Message.Error() + ": " + e.Obj
}

var ErrTypeHasNoFields = errors.New("type has no fields")
var ErrTypeMissingName = errors.New("missing name in type")
var ErrModelHasNoFields = errors.New("model has no fields")
var ErrModelMissingName = errors.New("missing name in model")
var ErrFieldMissingName = errors.New("missing field name in model")
var ErrUndefinedType = errors.New("type not defined")

func (c *Config) Sanitize() error {
	if len(c.Models) == 0 {
		return ErrNoModels
	}
	if len(c.Types) == 0 {
		log.Printf("no types defined") // HL7 types may not be defined
	}

	for i, t := range c.Types {
		if len(t.Fields) == 0 {
			return &CustomConfigError{ErrTypeHasNoFields, t.Name}
		}

		if t.Name == "" {
			return ErrTypeMissingName
		}

		for j, field := range t.Fields {
			c.validateImports(field.Type)

			if field.Tag == "" {
				c.Types[i].Fields[j].Tag = strconv.Itoa(j + 1)
			}

			c.Types[i].Fields[j].Name = capitalize(field.Name)

			if contains([]string{"uuid.UUID", "time.Time"}, mapType(field.Type)) {
				c.Types[i].needCustomUnmarshal = true
			}
		}

		c.Types[i].Name = capitalize(t.Name)
	}

	for i, model := range c.Models {
		if len(model.Fields) == 0 {
			return &CustomConfigError{ErrModelHasNoFields, model.Name}
		}
		if model.Name == "" {
			return ErrModelMissingName
		}

		c.Models[i].Name = capitalize(model.Name)
		for j, field := range model.Fields {
			if field.Name == "" {
				return &CustomConfigError{ErrFieldMissingName, model.Name}
			}

			if err := c.validateTypes(field.Type); err != nil {
				return err
			}

			c.validateImports(mapType(field.Type))

			c.Models[i].Fields[j].Name = capitalize(field.Name)
		}
	}

	return nil
}

func (c *Config) validateTypes(fieldType string) error {
	if contains([]string{"string", "int", "bool", "float", "time.Time", "uuid.UUID"}, mapType(fieldType)) {
		return nil
	}

	for _, t := range c.Types {
		if t.Name == fieldType {
			return nil
		}
	}

	return &CustomConfigError{ErrUndefinedType, fieldType}
}

func (c *Config) validateImports(fieldType string) {
	mapped := mapType(fieldType)
	if mapped == "time.Time" {
		if !contains(c.Meta.Imports, "time") {
			c.Meta.Imports = append(c.Meta.Imports, "time")
		}
		return
	}
}

func mapType(yamlType string) string {
	switch yamlType {
	case "":
		return "string"
	case "string":
		return "string"
	case "int":
		return "int"
	case "timestamp", "date":
		return "time.Time"
	case "uuid":
		return "uuid.UUID"
	default:
		return yamlType
	}
}
