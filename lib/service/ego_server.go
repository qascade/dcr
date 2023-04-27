package service

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func RunEgoServer() {
	// testCmds()
	fmt.Print("ego-server started on port 8080")

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/run", RunHandler)
	// http.Handle("/", http.FileServer(http.Dir("static")))

	http.ListenAndServe(":8080", nil)
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
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

	output := runTranformationGoCode()

	fmt.Fprintf(w, output)

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
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

func runCommands(commands []*exec.Cmd, dir string) (string, error) {
	// Change the working directory to the specified path
	err := os.Chdir(dir)
	if err != nil {
		return "", err
	}

	// Execute the commands in the specified directory
	var output string
	for _, cmd := range commands {
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		// if i == len(commands)-1 {
		output = string(out)
		// }
	}
	return output, nil
}

func runTranformationGoCode() string {

	fmt.Print("running tranformation..")

	dir := "../go_tranformations/test.zip"

	cmd_tidy := exec.Command("go", "mod", "tidy")
	cmd_build := exec.Command("ego-go", "build", "main.go")
	cmd_sign := exec.Command("ego", "sign", "main")
	cmd_generate_enclave := exec.Command("cp", "./enclave.json", "../../../go_tranformations/test.zip/", "-f")

	cmd_set_to_simulation := exec.Command("export", "OE_SIMULATION=1")
	cmd_run := exec.Command("OE_SIMULATION=1", "ego", "run", "main")

	// cmd_run := exec.Command("go", "run", "main.go")
	commands := []*exec.Cmd{cmd_tidy, cmd_build, cmd_sign, cmd_generate_enclave, cmd_set_to_simulation, cmd_run}

	output, err := runCommands(commands, dir)
	if err != nil {
		fmt.Println("Error running commands:", err)
		return ""
	}
	fmt.Println()
	fmt.Println("Output:", output)
	return output
}
