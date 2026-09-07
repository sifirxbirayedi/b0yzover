package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var version string

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Sürüm bilgisini yazdır",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("b0yzover sürüm: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
