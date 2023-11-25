package cmd

import (
	"finalproject/app/api"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use: "api",
	Run: func(_ *cobra.Command, _ []string) {
		api := api.New()
		api.Start()
	},
}
