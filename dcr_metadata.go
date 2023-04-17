package main

import (
	"os"
)

var (
	DCR_VERSION       string = "0.0.1"
	DCR_ENV_FILE_NAME string = ".env"
	DCR_DIR_NAME      string = ".dcr"
)

var (
	HOME_DIR          string
	DCR_DIR_FILE_PATH string
)

func init() {
	var err error
	HOME_DIR, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DCR_DIR_FILE_PATH = HOME_DIR + DCR_DIR_NAME
}
