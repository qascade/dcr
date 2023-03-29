/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	collab "github.com/qascade/dcr/collaboration"
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "Upload contracts",
	Short: "Upload contracts to the specified github repository",
	Long: `A cli tool to upload the contracts. The required parameters are the path where the contract is present
			and the repository to which it needs to be uploaded.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		contractPath := cmd.Flag("path").Value.String()
		collabPkg, err := collab.NewCollaborationPkg(contractPath)
		if err != nil {
			return err
		}
		collaboration := collab.Collaboration(collabPkg)
		err = collaboration.UploadToRepo(contractPath)
		if err != nil {
			return err
		}
		return nil
	},
}

var repoName string

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&contractPath, "path", "f", "", "Path to the contract file")
	uploadCmd.Flags().StringVarP(&repoName, "name", "f", "", "Name of Repository to upload the contract to")

}
