package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/change"
	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/project"
	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/proto"
	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/run"
	"github.com/go-kratos/kratos/cmd/kratos/v3/internal/upgrade"
)

var rootCmd = &cobra.Command{
	Use:           "kratos",
	Short:         "Kratos: An elegant toolkit for Go microservices.",
	Long:          `Kratos: An elegant toolkit for Go microservices.`,
	SilenceErrors: true,
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the kratos version information",
	Long:  "Show the kratos version information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(release)
	},
}

func init() {
	rootCmd.Version = release
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
