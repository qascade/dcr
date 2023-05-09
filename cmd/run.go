package cmd

import (
	"fmt"

	"github.com/qascade/dcr/lib/service"
	"github.com/spf13/cobra"
)

var (
	pkgPath string
)

// runCmd represents the run command
// dcr run -p <pkgpath>
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Transformation and Destination",
	Long:  `CLI Command to run mentioned transformation and destination`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgPath := cmd.Flag("pkgpath").Value.String()
		service, err := service.NewService(pkgPath)
		if err != nil {
			err = fmt.Errorf("err creating new service with package path: %s", pkgPath)
			return err
		}
		err = service.Run()
		if err != nil {
			err := fmt.Errorf("err running service: %s", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("pkgpath", "p", "", "reference of the destination")
}
