package cmd

import (
        "github.com/spf13/cobra"
        "fmt"
        "github.com/jacobtomlinson/krontab/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
        Use:   "version",
        Short: "Print the version number of krontab",
        Long:  `All software has versions. This is krontab`,
        Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Build Date:", version.BuildDate)
                fmt.Println("Git Commit:", version.GitCommit)
                fmt.Println("Version:", version.Version)
                fmt.Println("Go Version:", version.GoVersion)
                fmt.Println("OS / Arch:", version.OsArch)
        },
}

func init() {
        rootCmd.AddCommand(versionCmd)
}
