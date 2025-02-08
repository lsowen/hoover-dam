package cmd

import (
	"fmt"
	"github.com/lsowen/hoover-dam/pkg/api"
	"github.com/spf13/cobra"
	"net/http"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run hoover-dam",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("loading run command config: %w", err)
		}

		r, err := api.Serve(cmd.Context(), *cfg)
		if err != nil {
			return fmt.Errorf("serving api: %w", err)
		}

		err = http.ListenAndServe(":8080", r)
		if err != nil {
			return fmt.Errorf("listening and serving api on port 8080: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
