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
	Long: `Krontab is a crontab replacement for kubernetes.

You can use it to create cron jobs on a kubernetes cluster in a familiar crontab format.
Krontab works by allowing you to create job templates which are used in kubernetes. Then create
specific cron jobs using the crontab. Example:

# Crontab example

# template: default
0 1 * * * echo hello  # name: test
`,
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

	rootCmd.Flags().BoolVarP(&editCrontab, "edit-crontab", "e", false, "Edit the crontab")
	rootCmd.Flags().BoolVarP(&listCrontab, "list-crontab", "l", false, "List the crontab")
}
