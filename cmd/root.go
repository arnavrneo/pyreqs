/*
Copyright © 2023 ARNAV <r.arnav@icloud.com>
*/
package cmd

import (
	"os"
	"pyreqs/cmd/create"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pyreqs",
	Short: "pyreqs: Create python dependency file easily",
	Long: `pyreqs: A cli program designed to speed-up your python projects by auto-generating python dependencies (requirements.txt) file for you. 
Call pyreqs with the available commands to see the magic in work!
			`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandsPalette() {
	rootCmd.AddCommand(create.CreateCmd)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pyreqs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	addSubcommandsPalette()
}
