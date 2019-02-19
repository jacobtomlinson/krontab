package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jacobtomlinson/krontab/template"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a krontab resource",
	Long:  `Create a krontab resource`,
}

var createTemplaceCmd = &cobra.Command{
	Use:   "template",
	Short: "Create a template resource",
	Long:  `Create a template resource`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template.CreateTemplate(args[0])
	},
}

func init() {
	createCmd.AddCommand(createTemplaceCmd)
	rootCmd.AddCommand(createCmd)
}
