package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jacobtomlinson/krontab/crontab"
)

var jobCommand bool
var jobTemplate string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a krontab job",
	Long:  `Run a krontab job`,
}

var runJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Run a krontab job manually",
	Long:  `Run a krontab job manually`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if jobCommand {
			_, err := crontab.RunJob(args, jobTemplate)
			if err != nil {
				os.Exit(1)
			}
		} else {
			_, err := crontab.RunCronJob(args[0])
			if err != nil {
				os.Exit(1)
			}
		}
	},
}

func init() {
	runCmd.AddCommand(runJobCmd)
	runJobCmd.Flags().BoolVarP(&jobCommand, "command", "c", false, "Specify a one shot command to run")
	runJobCmd.Flags().StringVar(&jobTemplate, "template", "default", "Specify a template to use for one shot commands")
	rootCmd.AddCommand(runCmd)
}
