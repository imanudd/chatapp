package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "application",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

func Run() {
	rootCmd.AddCommand(apiCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
