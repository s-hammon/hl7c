package config

import "fmt"

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
