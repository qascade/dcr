
package cmd

import (
	"fmt"
	"embed"
    "io/fs"
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
		dublicateSampleFolder()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

//go:embed source_folder/*
var sampleFolder embed.FS

func dublicateSampleFolder(){
	destDir := "temp2"
	sourceFolderPath := "temp"
	//samples\hw_collaboration

	err = copyFolder(sourceFolderPath, destDir)    
}

func copyFolder(sourceFolderPath, destinationFolderPath string) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}

	// Construct the absolute paths to the source and destination folders
	sourceAbsPath := filepath.Join(cwd, sourceFolderPath)
	destAbsPath := filepath.Join(cwd, destinationFolderPath)

	// Open the source folder
	source, err := embed.Embed(sourceAbsPath)
	if err != nil {
		return fmt.Errorf("could not open source folder: %w", err)
	}

	// Create the destination folder if it doesn't exist
	if _, err := os.Stat(destAbsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(destAbsPath, 0755); err != nil {
			return fmt.Errorf("could not create destination folder: %w", err)
		}
	}

	// Copy files from the source to the destination folder
	return filepath.Walk(sourceAbsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not access file %q: %w", path, err)
		}

		// Skip directories and hidden files
		if info.IsDir() || filepath.Base(path)[0] == '.' {
			return nil
		}

		// Open the source file
		file, err := source.Open(path)
		if err != nil {
			return fmt.Errorf("could not open source file %q: %w", path, err)
		}
		defer file.Close()

		// Create the destination file
		destPath := filepath.Join(destAbsPath, path[len(sourceAbsPath):])
		destDir := filepath.Dir(destPath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Errorf("could not create destination directory %q: %w", destDir, err)
			}
		}
		dest, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("could not create destination file %q: %w", destPath, err)
		}
		defer dest.Close()

		// Copy the file contents
		if _, err := io.Copy(dest, file); err != nil {
			return fmt.Errorf("could not copy file contents from %q to %q: %w", path, destPath, err)
		}

		return nil
	})
}