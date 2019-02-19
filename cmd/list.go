package cmd

import (
        "github.com/spf13/cobra"
		"fmt"

		"github.com/jacobtomlinson/krontab/template"
		"github.com/jacobtomlinson/krontab/crontab"
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

func init() {
	listCmd.AddCommand(listTemplaceCmd)
	listCmd.AddCommand(listCrontabCmd)
	rootCmd.AddCommand(listCmd)
}
