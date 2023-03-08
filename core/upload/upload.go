package main

import (
	"context"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/joho/godotenv"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading environment variables file")
		os.Exit(1)
	}

	ctx := context.Background()
	token := os.Getenv("token")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoName := "examplei"
	repoDescription := "This is a new repo created via the GitHub API"
	repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		Name:        &repoName,
		Description: &repoDescription,
	})
	if err != nil {
		fmt.Printf("Failed to create repo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created new repo: %v\n", *repo.HTMLURL)

	owner := os.Getenv("owner")
    repoT := "examplei"
    filePath := "example.txt"

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fileOptions := &github.RepositoryContentFileOptions{
        Message:   github.String("Add example file"),
        Content:   data,
    }

    _, _, err = client.Repositories.CreateFile(
        context.Background(),
        owner,
        repoT,
        filePath,
        fileOptions,
    )
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("File uploaded successfully")

	

}