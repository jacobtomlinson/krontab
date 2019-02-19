package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jacobtomlinson/krontab/template"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a krontab resource",
	Long:  `Delete a krontab resource`,
}

var deleteTemplaceCmd = &cobra.Command{
	Use:   "template",
	Short: "Delete a template resource",
	Long:  `Delete a template resource`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template.DeleteTemplate(args[0])
	},
}

func init() {
	deleteCmd.AddCommand(deleteTemplaceCmd)
	rootCmd.AddCommand(deleteCmd)
}
