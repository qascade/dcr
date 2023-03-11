package cmd

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
// this cmd will generate colaboration sample config files in folder named "helloworld_configs" in same directory
// configs files name will be : "collaborator1_tables.yaml", "collaborator2_tables.yaml", "contract.yaml"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init called")
		generateCollabFiles()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

//go:embed collabsamples/hw_collaboration/*
var f embed.FS

// this func, uses embed package and provides files inside [samples/hw_collaboration/*]
// during run time.
// Files are then saved in a new directory
func generateCollabFiles() {
	contractFile, err := f.ReadFile("collabsamples/hw_collaboration/contract.yaml")
	if err != nil {
		log.Fatal(err)
	}
	collaborator1_tables, err := f.ReadFile("collabsamples/hw_collaboration/collaborator1_tables.yaml")
	if err != nil {
		log.Fatal(err)
	}
	collaborator2_tables, err := f.ReadFile("collabsamples/hw_collaboration/collaborator2_tables.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// print(string(data))

	outputDir := "collab_pkg"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	writeFilesToGivenOutputDir(outputDir, contractFile, collaborator1_tables, collaborator2_tables)

}

// writing files to new directory, files data is fetched at runtime using embed package
func writeFilesToGivenOutputDir(outputDir string, contractFile []byte, collaborator1_tables []byte, collaborator2_tables []byte) {
	err := ioutil.WriteFile(filepath.Join(outputDir, "contract.yaml"), contractFile, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(outputDir, "collaborator1_tables.yaml"), collaborator1_tables, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(outputDir, "collaborator2_tables.yaml"), collaborator2_tables, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File saved successfully!")
}
