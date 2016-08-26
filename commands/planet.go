package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

//CfgFile is exported
var CfgFile string

func init() {

	cobra.OnInitialize(initConfig)
	RootCMD.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/dagobah/config.yaml)")
}

func initConfig() {

	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("etc/dagobah/")
	viper.AddConfigPath("etc/dagobah/")
	viper.ReadInConfig()
}
