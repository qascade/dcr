package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v37/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading environment variables file")
		os.Exit(1)
	}

	owner := os.Getenv("owner")
	repo := "examplei"
	token := os.Getenv("token")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, er := client.Repositories.Delete(ctx, owner, repo)
	if er != nil {
		fmt.Println("Error deleting repository:", er)
		os.Exit(1)
	}

	fmt.Println("Repository deleted successfully")
}
