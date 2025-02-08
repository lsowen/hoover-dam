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

var rootCmd = &cobra.Command{
	Use:   "hoover-dam",
	Short: "hoover-dam is an open source authorization server for lakefs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("executing command: %v\n", err)
		os.Exit(1)
	}
}

var initOnce sync.Once

func newConfig() (*config.Config, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("initializing new config: %w", err)
	}

	return cfg, nil
}

func loadConfig() (*config.Config, error) {
	initOnce.Do(initConfig)
	cfg, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	return cfg, nil
}

func initConfig() {
	// Use experimental feature in 1.20 alpha https://github.com/spf13/viper/issues/1851
	viper.SetOptions(viper.ExperimentalBindStruct())
	viper.SetEnvPrefix("HOOVERDAM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
