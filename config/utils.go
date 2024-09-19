package config

import (
	"fmt"
	"strings"
)

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func toGoTag(s string) string {
	return fmt.Sprintf("`json:\"%s\"`", s)
}

func createBaseStruct() string {
	var sb strings.Builder
	sb.WriteString("type Base struct {\n")
	sb.WriteString("\tID        uuid.UUID\n")
	sb.WriteString("\tCreatedAt time.Time\n")
	sb.WriteString("\tUpdatedAt time.Time\n")
	sb.WriteString("}\n\n")

	return sb.String()
}
