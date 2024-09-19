package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hl7c",
	Short: "hl7c is a CLI tool for generating HL7 models",
	Long: `hl7c is a CLI tool for generating HL7 models from a config file.
	
HL7 is a messaging standard that enables clinical applications to exchange data.
With hl7c, you can generate Go models from a config file which will read incoming
messages (in JSON format) and unmarshal them into Go structs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
