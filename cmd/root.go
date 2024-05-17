// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/wherehouse/pkg/configurator"
)

const rootCommandLongDesc = "wherehouse is a tracking application for personal items.\n" +
	"It stores a digital record of items with options to use, delete, loan, and borrow."

// application build information set by the linker.
var (
	Version string
	Date    string
)

var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := GetRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetRootCommand returns the root command for the CLI.
func GetRootCommand() *cobra.Command {
	var cfgFile string
	if rootCmd == nil {
		rootCmd = &cobra.Command{
			Use:     "wherehouse",
			Version: fmt.Sprintf("%s\n%s", Version, Date),
			Short:   "wherehouse is an inventory tracking application",
			Long:    rootCommandLongDesc,
			Run: func(cmd *cobra.Command, args []string) {
				cfg, err := configurator.New(configurator.WithFile(cfgFile))
				if err != nil {
					cmd.PrintErr(err)
				}

				cmd.Println("config file:", cfg.GetConfigFileUsed())
			},
		}
	}

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "configuration file")

	// rootCmd.AddCommand(GetAddCmd())
	// rootCmd.AddCommand(GetInfoCmd())
	// rootCmd.AddCommand(GetReportCmd())
	// rootCmd.AddCommand(GetRemoveCmd())
	// rootCmd.AddCommand(GetUpdateCmd())

	return rootCmd
}
