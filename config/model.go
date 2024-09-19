package config

import (
	"fmt"
	"strings"
)

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
			customFields[field.Name] = fmt.Sprintf("time.Parse(\"20060102150405\", aux.%s)", field.Name)
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
