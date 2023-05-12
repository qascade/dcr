package service

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/qascade/dcr/lib/utils"
)

func generateScriptAndRun() string {

	go_tranformation_path := "../go_tranformations/test.zip"
	utils.ChangeDirTo(go_tranformation_path)

	// Define the script content as a string
	script := generateScript()

	// Write the script to a file
	err := os.WriteFile("myscript.sh", []byte(script), 0755)
	if err != nil {
		fmt.Println("Error writing script:", err)
		return "error writing script"
	}

	// Make the script executable
	cmd := exec.Command("chmod", "+x", "myscript.sh")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error making script executable:", err)
		return "error making script executable"
	}

	// Run the script and capture the output
	out := runScript("./myscript.sh")

	// Print the output
	fmt.Println(string(out))

	utils.DeleteTempFolder("../go_tranformations/test.zip")

	return string(out)
}

func generateScript() string {
	script :=
		`
#!/bin/bash

echo "Starting script to run go-tranformation.."

echo >> "running : [go mod tidy] to load dependencies of tranformation"
go mod tidy

echo >> "running : [ego-go build main.go] 
ego-go build main.go


ego sign main

echo >> "running : [modifying enclave.json config file] to include datafiles"
cp ../../ego_server/enclave.json ./ -f
ego sign main

echo >> "running : [OE_SIMULATION=1 ego run main] to set mode to SIMULATION"
OE_SIMULATION=1 ego run main
`
	return script
}

func runScript(path string) string {
	cmd, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	return output
}

func testScript() string {
	go_tranformation_path := "../go_tranformations/test.zip"
	utils.ChangeDirTo(go_tranformation_path)

	// Run the script and capture the output
	out := runScript("./myscript.sh")
	return string(out)
}
