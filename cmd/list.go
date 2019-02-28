package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jacobtomlinson/krontab/crontab"
	"github.com/jacobtomlinson/krontab/template"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List krontab resources",
	Long:  `List krontab resources`,
}

var listTemplaceCmd = &cobra.Command{
	Use:   "template",
	Short: "List a template resource",
	Long:  `List a template resource`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, element := range template.ListTemplates() {
			fmt.Println(element)
		}
	},
}

var listCrontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: "list the crontab",
	Long:  `list the crontab`,
	Run: func(cmd *cobra.Command, args []string) {
		crontab.ListCrontab()
	},
}

var listRunningCmd = &cobra.Command{
	Use:   "running",
	Short: "list the running jobs",
	Long:  `list the running jobs`,
	Run: func(cmd *cobra.Command, args []string) {
		jobs, _ := crontab.ListRunning()
		if len(jobs) > 0 {
			for _, job := range jobs {
				fmt.Println(job)
			}
		} else {
			fmt.Println("No running jobs.")
		}
	},
}

func init() {
	listCmd.AddCommand(listRunningCmd)
	listCmd.AddCommand(listTemplaceCmd)
	listCmd.AddCommand(listCrontabCmd)
	rootCmd.AddCommand(listCmd)
}
