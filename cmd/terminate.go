/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var terminateCmd = &cobra.Command{
	Use:   "Terminate the contract between collaborators",
	Short: "Terminate and delete the contract between the collaborators",
	Long: `A cli tool to terminate the contracts. The required parameters is the name of the Repository in 
			which the collaborator had uploaded the contract. `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	rootCmd.AddCommand(terminateCmd)
	terminateCmd.Flags().StringVarP(&repoName, "name", "f", "", "Name of Repository to upload the contract to")

}
