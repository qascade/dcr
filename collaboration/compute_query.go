package collaboration

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/snowflakedb/gosnowflake"
)

func (c *CollaborationPackage) ExecuteWarehouse(pathToSQLfile string) error {

	file, err := os.Open(pathToSQLfile)

	if err != nil {
		fmt.Println("Unable to access SQL file.")
		return err
	}

	defer file.Close()

	sqlQuery, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Unable to read SQL file")
		return err
	}

	err = godotenv.Load("../.env")
	if err != nil {
		err = errors.New("error loading environment variables file")
		return err
	}

	warehouseConfig := gosnowflake.Config{
		Account:   os.Getenv("Account"),
		User:      os.Getenv("User"),
		Password:  os.Getenv("Password"),
		Database:  os.Getenv("Database"),
		Schema:    os.Getenv("Schema"),
		Warehouse: os.Getenv("Warehouse"),
	}

	dsnString, err := gosnowflake.DSN(&warehouseConfig)

	if err != nil {
		fmt.Println("Error dumping config file to DSN format")
	}

	db, err := sql.Open("snowflake", dsnString)
	if err != nil {
		fmt.Println("Error accessing the snowflake warehouse")
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(sqlQuery))
	if err != nil {
		fmt.Println("Unabe to execute sql query: ", err)
		return err
	}

	log.Println("Query executed successfully.")
	return nil
}
