package utils

import (
	"fmt"
	"bytes"
	"os/exec"
	"os"
	"io"

	"gopkg.in/yaml.v3"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func UnmarshalStrict(in []byte, out interface{}) (err error) {
	knownFieldsDecoder := yaml.NewDecoder(bytes.NewReader(in))
	knownFieldsDecoder.KnownFields(true)
	return knownFieldsDecoder.Decode(out)
}

func RunCmd(cmd *exec.Cmd) (string,error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to capture output for command:%s", err)
	}
	log.Infof("Output of cmd %s: %s", cmd.String(), output)
	return string(output), err
}

// Give full path to the new file including the new name.
func CopyFile(dest string, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("not able to open file:%s ", src)
		log.Error(err)
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		err = fmt.Errorf("not able to create file:%s ", dest)
		log.Error(err)
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		err = fmt.Errorf("unable to copy %s to %s", src, dest)
		log.Error(err)
		return err
	}
	log.Infof("copyfile from %s to %s", src, dest)
	return nil
}

func Remove(filePath string) error {
	if !FileExists(filePath) {
		return nil
	}
	err := os.Remove(filePath)
	if err != nil {
		err = fmt.Errorf("unable to remove file: %s", filePath)
		log.Error(err)
		return err
	}
	return nil
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
    if os.IsNotExist(err) {
       return false
    }
    return !info.IsDir()	
}

func WriteStringToFile(filePath string, content string) error {
	newFile, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("unable to create file %s", filePath)
		log.Error(err)
		return err
	}

	_, err = io.WriteString(newFile, content)
	if err != nil {
		err = fmt.Errorf("unable to write to file %s", filePath)
		log.Error(err)
		return err
	}
	log.Infof("Writing to file %s", filePath)
	return nil
}