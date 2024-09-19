package config

import (
	"fmt"
	"strings"
)

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
		if field.Type == "time.Time" {
			customFields[field.Name] = fmt.Sprintf("time.Parse(\"20060102150405\", aux.%s)", field.Name)
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

func mapType(yamlType string) string {
	switch yamlType {
	case "":
		return "string"
	case "string":
		return "string"
	case "int":
		return "int"
	case "timestamp":
		return "time.Time"
	case "uuid":
		return "uuid.UUID"
	default:
		return yamlType
	}
}
