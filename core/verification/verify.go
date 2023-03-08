package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v35/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

func extract(link string) string {
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("token")},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(tokenClient)

	owner := os.Getenv("owner")
	repo := link

	filePath := "contract.yaml"

	fileContents, _, _, err := client.Repositories.GetContents(
		ctx, owner, repo, filePath, nil,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	content, err := fileContents.GetContent()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(content)
	return content
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading environment variables file")
		os.Exit(1)
	}

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

	var git_repos [100]string

	var n = len(yamlConfig.Collaborator)

	for i := 0; i < n; i++ {
		git_repos[i] = yamlConfig.Collaborator[i].ContractRepo
	}

	var ini string

	for i := 0; i < n; i++ {

		if i == 0 {
			ini = extract(git_repos[i])
			continue
		}
		temp := extract(git_repos[i])

		if ini != temp {
			fmt.Println("Content mismatch error")
			os.Exit(1)
		}

	}

	fmt.Println("Contracts verified successfully")

}
