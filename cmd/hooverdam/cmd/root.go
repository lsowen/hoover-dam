package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/lsowen/hoover-dam/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hoover-dam",
	Short: "hoover-dam is an open source authorization server for lakefs",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var initOnce sync.Once

func newConfig() (*config.Config, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadConfig() *config.Config {
	initOnce.Do(initConfig)
	cfg, err := newConfig()
	if err != nil {
		fmt.Println("Failed to load config file", err)
		os.Exit(1)
	}
	return cfg
}

func initConfig() {

	// Use experimental feature in 1.20 alpha https://github.com/spf13/viper/issues/1851
	viper.SetOptions(viper.ExperimentalBindStruct())

	viper.SetEnvPrefix("HOOVERDAM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // support nested config
	// read in environment variables
	viper.AutomaticEnv()

}
