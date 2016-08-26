package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCMD Exported Variable
var RootCMD = &cobra.Command{

	Use:   "Dagobah",
	Short: `Dagobah is an awesome planet style RSS aggregator`,
	Long:  `Dagobah provides planet style RSS aggregation. It is inspired by python planet. It has a simple YAML configuration and provides it's own webserver.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dagobah Runs")
	},
}

//Execute will run the cobra commands.
func Execute() {

	err := RootCMD.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
