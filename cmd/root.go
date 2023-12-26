/*
Copyright © 2023 ARNAV <r.arnav@icloud.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"pyreqs/cmd/clean"
	"pyreqs/cmd/create"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pyreqs",
	Short: "pyreqs: Create python dependencies file easily",
	Long: `Seamlessly generate python dependencies (requirements.txt) file. 
Call pyreqs with the available commands to see the magic done!
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandsPalette() {
	rootCmd.AddCommand(create.Cmd)
	rootCmd.AddCommand(clean.Cmd)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("dirPath", "d", "./", "directory to .py files")
	rootCmd.PersistentFlags().StringP("venvPath", "v", " ", "directory to venv (virtual env)")
	rootCmd.PersistentFlags().StringP("ignore", "i", " ", "ignore specific directories (each seperated by comma)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	addSubcommandsPalette()
}
