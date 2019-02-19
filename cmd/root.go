package cmd

import (
	"fmt"
	"os"

	"github.com/jacobtomlinson/krontab/crontab"
	"github.com/spf13/cobra"
)

var cfgFile string
var editCrontab bool
var listCrontab bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krontab",
	Short: "A crontab replacement for kubernetes",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if editCrontab {
			crontab.EditCrontab()
		} else if listCrontab {
			crontab.ListCrontab()
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	// TODO Flags are not working
	rootCmd.Flags().BoolVarP(&editCrontab, "edit-crontab", "e", false, "Edit the crontab")
	rootCmd.Flags().BoolVarP(&listCrontab, "list-crontab", "l", false, "List the crontab")
}
