/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/qascade/dcr/lib/service"
	"github.com/spf13/cobra"
)

// initServerCmd represents the initServer command
var initServerCmd = &cobra.Command{
	Use:   "initServer",
	Short: "Starts confidential clean room server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.RunEgoServer()
	},
}

func init() {
	rootCmd.AddCommand(initServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
