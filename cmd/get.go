package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jacobtomlinson/krontab/crontab"
	"github.com/jacobtomlinson/krontab/template"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get krontab resources",
	Long:  `Get krontab resources`,
}

var getTemplaceCmd = &cobra.Command{
	Use:   "template",
	Short: "Get a template resource",
	Long:  `Get a template resource`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, _ := template.GetTemplate(args[0])
		fmt.Println(strings.TrimSpace(tmpl))
	},
}

var getCrontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: "Get the crontab",
	Long:  `Get the crontab`,
	Run: func(cmd *cobra.Command, args []string) {
		crontab.ListCrontab()
	},
}

func init() {
	getCmd.AddCommand(getTemplaceCmd)
	getCmd.AddCommand(getCrontabCmd)
	rootCmd.AddCommand(getCmd)
}
