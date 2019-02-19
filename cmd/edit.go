package cmd

import (
    "github.com/spf13/cobra"

    "github.com/jacobtomlinson/krontab/crontab"
    "github.com/jacobtomlinson/krontab/template"
)

var editCmd = &cobra.Command{
    Use:   "edit",
    Short: "Edit a krontab resource",
    Long:  `Edit a krontab resource`,
}

var editTemplaceCmd = &cobra.Command{
	Use:   "template",
	Short: "Edit a template resource",
	Long:  `Edit a template resource`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Arg checking in template editing
		template.EditTemplate(args[0])
	},
}

var editCrontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: "Edit the crontab",
	Long:  `Edit the crontab`,
	Run: func(cmd *cobra.Command, args []string) {
        crontab.EditCrontab()
	},
}

func init() {
	editCmd.AddCommand(editTemplaceCmd)
	editCmd.AddCommand(editCrontabCmd)
	rootCmd.AddCommand(editCmd)
}
