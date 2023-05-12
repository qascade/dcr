package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SendRequestToServer(folderPath string, ref string) error {
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

	// Create a multipart form with the ZIP file and the string value
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the string value as a form field
	err = multipartWriter.WriteField("refid", ref)
	if err != nil {
		return err
	}

	// Add the ZIP file as a form field
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
	req, err := http.NewRequest("POST", "https://bb53-117-98-96-74.in.ngrok.io/run", &requestBody)
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

func AskForResult(addressRef string) (string, error) {
	// Set the URL with the refid parameter
	// refid := "123456"
	url := fmt.Sprintf("https://bb53-117-98-96-74.in.ngrok.io/download?refid=%s", addressRef)

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)

	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)

	}

	// Print the response body
	fmt.Println(string(body))
	return "", nil
}
