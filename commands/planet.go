package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCMD is a cobra command that will be executed main
var RootCMD = &cobra.Command{

	Use:   "...",
	Short: `...`,
	Long:  `...`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {

	fmt.Println(viper.Get("feeds"))
	fmt.Println(viper.GetString("appname"))

}

//Execute will run the RootCMD command  but just before, it will add more commands.
func Execute() {

	addCommands()

	err := RootCMD.Execute()
	if err != nil {
		fmt.Println(err)
		fmt.Println("something is wrong ")
		os.Exit(-1)
	}
}

func addCommands() {

	RootCMD.AddCommand(fetchCmd)
}

//CfgFile is the name of the variable that will stores the configuration.
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
	viper.AddConfigPath("$HOME/.dagobah/")
	viper.ReadInConfig()
}
