package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func PrepareToSend() {
	folderPath := "./go_app_templates/go_app"
	err := SendRequestToServer(folderPath)
	if err != nil {
		log.Fatal(err)
	}
}

func SendRequestToServer(folderPath string) error {
	// Create a buffer to store the ZIP file contents
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Walk through the folder and add all files to the ZIP
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			relPath, err := filepath.Rel(folderPath, path)
			if err != nil {
				return err
			}
			file, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}
			inFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer inFile.Close()
			_, err = io.Copy(file, inFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Close the ZIP writer
	err = zipWriter.Close()
	if err != nil {
		return err
	}

	// Create a multipart form with the ZIP file
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	zipFile, err := multipartWriter.CreateFormFile("file", "test.zip")
	if err != nil {
		return err
	}
	_, err = io.Copy(zipFile, &buf)
	if err != nil {
		return err
	}
	err = multipartWriter.Close()
	if err != nil {
		return err
	}

	// Create a POST request to the /run endpoint of the API
	req, err := http.NewRequest("POST", "http://localhost:8080/run", &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request and check the response status code
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}
