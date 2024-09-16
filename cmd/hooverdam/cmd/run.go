package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lsowen/hoover-dam/pkg/api"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run hoover-dam",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig()

		r, err := api.Serve(cmd.Context(), *cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = http.ListenAndServe(":8080", r)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
