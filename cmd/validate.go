package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "This command will parse and validate the contract package with the given path",
	Long: `
This command will parse and validate the contract package with the given path.
This command will check for name mismathces or table duplications or essential fields missing and return error if any.
	`,
	RunE: Validate,
}

func Validate(cmd *cobra.Command, args []string) error {
	// Takes the contract package path as input and parse yaml
	contractPath = cmd.Flag("path").Value.String()
	fmt.Println(contractPath)
	return nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVarP(&contractPath, "path", "f", "", "Path to the contract package")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
