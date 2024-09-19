package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/s-hammon/hl7c/config"
	"github.com/spf13/cobra"
)

var pkgPath string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "'generate' will create a model.go file in ./internal/{package} directory",
	Long: `generate will create a model.go file in ./internal/{package} directory:
	
'<cmd> generate -f config.yaml -p internal/objects'
	
This will create a model.go file at internal/objects in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("file")
		abspath, err := filepath.Abs(configFile)
		if err != nil {
			fmt.Println("error reading config file: ", err)
			os.Exit(1)
		}

		pkg, _ := cmd.Flags().GetString("package")
		models := genModels(abspath, pkg)

		pkgPath, err = saveFile(pkg, models)
		if err != nil {
			fmt.Println("error saving file: ", err)
			os.Exit(1)
		}

		shiageru()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("file", "f", "model_config.yaml", "config file to generate models from")
	generateCmd.PersistentFlags().StringP("package", "p", "objects", "package to generate models in")
}

func shiageru() {
	command := exec.Command("go", "fmt", pkgPath)
	if err := command.Run(); err != nil {
		fmt.Println("error formatting models.go: ", err.Error())
		os.Exit(1)
	}

	command = exec.Command("go", "get", "./...")
	if err := command.Run(); err != nil {
		fmt.Println("error getting dependencies: ", err.Error())
		os.Exit(1)
	}
}

func genModels(abspath, pkg string) string {
	data, err := os.ReadFile(abspath)
	if err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.New(data, pkg)
	if err != nil {
		fmt.Println("error creating config: ", err)
		os.Exit(1)
	}

	return cfg.Compile()
}

func saveFile(pkg, data string) (string, error) {
	dir := fmt.Sprintf("internal/%s", pkg)
	if _, err := os.Stat(pkg); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}

	filePath := path.Join(dir, "model.go")
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}
