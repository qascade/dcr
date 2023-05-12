package service

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RunHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("POST request recieved at /run \n")
	fmt.Print("processing request.. \n")
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Print(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(32 << 20) // Max file size 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ZIP file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary directory to extract the files from the ZIP
	tmpDir, err := os.MkdirTemp("", "zip-extract")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	// Extract the files from the ZIP to the temporary directory
	zipReader, err := zip.NewReader(file, handler.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, f := range zipReader.File {
		fpath := filepath.Join(tmpDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			http.Error(w, "Invalid ZIP file", http.StatusBadRequest)
			return
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		inFile, err := f.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer inFile.Close()
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Save the extracted files to the server directory
	destDir := "../go_tranformations"
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.Rename(tmpDir, filepath.Join(destDir, handler.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message

	// output := runTranformationGoCode()
	output := generateScriptAndRun()

	fmt.Print("logging output..\n")
	fmt.Print(output)

	fmt.Fprintf(w, output)

}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<html>
				<body>
					<form action="/run" method="post" enctype="multipart/form-data">
						<input type="file" name="file"><br><br>
						<input type="submit" value="Run">
					</form>
				</body>
			</html>`
	fmt.Fprint(w, html)
}
