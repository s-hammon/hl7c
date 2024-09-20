package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/trace"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Do(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "hl7c", SilenceUsage: true}
	rootCmd.PersistentFlags().StringP("file", "f", "", "specify an alternative config file (default: model_config.yaml)")

	rootCmd.AddCommand(generateCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 1
		}
	}

	return 0
}

func getConfigPath(stderr io.Writer, f *pflag.Flag) (string, string) {
	if f != nil && f.Changed {
		file := f.Value.String()
		if file == "" {
			fmt.Fprintln(stderr, "error parsing config: file argument is empty")
			os.Exit(1)
		}
		abspath, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(stderr, "error parsing config: absolute path lookup failed: %s\n", err)
			os.Exit(1)
		}
		return filepath.Dir(abspath), filepath.Base(abspath)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing model_config.yaml: file does not exist")
			os.Exit(1)
		}
		return wd, ""
	}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate source code from a config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "generate").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		output, err := Generate(cmd.Context(), dir, name, &Options{
			Stderr: stderr,
		})
		if err != nil {
			os.Exit(1)
		}
		defer trace.StartRegion(cmd.Context(), "writeFile").End()
		for path, data := range output {
			if err := saveFile(path, data); err != nil {
				fmt.Fprintf(stderr, "%s: %s\n", path, err)
				os.Exit(1)
			}
		}

		shiageru()
		return nil
	},
}

func saveFile(dir, data string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	filePath := path.Join(dir, "model.go")
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}

func shiageru() {
	command := exec.Command("go", "fmt", "./...")
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
