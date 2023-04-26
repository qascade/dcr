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
// this cmd will generate colaboration sample config files in folder named "hw_collaboration" in same directory
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
		destPath := cmd.Flag("dest").Value.String()
		generateCollabFiles(destPath)
		return nil
	},
}

var collabDestPath string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&collabDestPath, "dest", "d", "", "Destination path to initialize collab pkg")
}

//go:embed collabsamples/hw_collaboration/*
var f embed.FS

// this func, uses embed package and provides files inside [samples/hw_collaboration/*]
// during run time.
// Files are then saved in a new directory
func generateCollabFiles(destPath string) error {
	contractFile, err := f.ReadFile("collabsamples/hw_collaboration/contract.yaml")
	if err != nil {
		log.Fatal(err)
		return err
	}
	collaborator1_tables, err := f.ReadFile("collabsamples/hw_collaboration/collaborator1_tables.yaml")
	if err != nil {
		log.Fatal(err)
		return err
	}
	collaborator2_tables, err := f.ReadFile("collabsamples/hw_collaboration/collaborator2_tables.yaml")
	if err != nil {
		log.Fatal(err)
		return err
	}
	// print(string(data))

	outputDir := destPath + "/collab_files"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return writeFilesToGivenOutputDir(outputDir, contractFile, collaborator1_tables, collaborator2_tables)
}

// writing files to new directory, files data is fetched at runtime using embed package
func writeFilesToGivenOutputDir(outputDir string, contractFile []byte, collaborator1_tables []byte, collaborator2_tables []byte) error {
	err := ioutil.WriteFile(filepath.Join(outputDir, "contract.yaml"), contractFile, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(outputDir, "collaborator1_tables.yaml"), collaborator1_tables, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(outputDir, "collaborator2_tables.yaml"), collaborator2_tables, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("File saved successfully! %s", outputDir)
	return nil
}
