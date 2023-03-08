package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type contract struct {
    Name         string `yaml:"name"`
    Version      string `yaml:"version"`
    Purpose      string `yaml:"Purpose"`
    Collaborator []struct {
        Name         string `yaml:"name"`
        ContractRepo string `yaml:"contract_repo"`
        UserAgents   []struct {
            Name  string `yaml:"name"`
            Email string `yaml:"email"`
        } `yaml:"user_agents"`
        Warehouse []struct {
            Name string `yaml:"name"`
        } `yaml:"warehouse"`
    } `yaml:"collaborators"`
    ComputeWarehouse string `yaml:"compute_warehouse"`
}


func main() {
    fmt.Println("Parsing YAML file")

    var fileName string
    flag.StringVar(&fileName, "f", "", "YAML file to parse.")
    flag.Parse()

    if fileName == "" {
        fmt.Println("Please provide yaml file by using -f option")
        return
    }

    yamlFile, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Printf("Error reading YAML file: %s\n", err)
        return
    }

    var yamlConfig contract
    err = yaml.Unmarshal(yamlFile, &yamlConfig)
    if err != nil {
        fmt.Printf("Error parsing YAML file: %s\n", err)
    } else {
		fmt.Printf("Parsed successfully\n")
	}

    fmt.Printf("Result: %v\n", yamlConfig)
}