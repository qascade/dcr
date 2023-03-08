package collaboration

import (
	"context"
	"fmt"
	"os"
	"errors"

	"github.com/google/go-github/v35/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)



func GetContractFromRepo(link string) (*GitRepoContent, error) {
	
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
		err = errors.New("error getting file contents from GitRepo")
		return nil, err
	}

	content, err := fileContents.GetContent()
	if err != nil {
		err = errors.New("error getting file content from Contract")
		return nil, err
	}

	fmt.Println(content)
}

func Verify(path string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading environment variables file")
	}

	var collabPkg CollaborationParser = &CollaborationPackage{}
	cSpec, _, err := collabPkg.Parse(path)
	if err != nil {
		return fmt.Errorf("error parsing contract.yaml file, %v", err)
	}


	var gitRepos [100]string

	var noOfCollaborators = len(cSpec.Collaborators)

	for i := 0; i < noOfCollaborators; i++ {
		gitRepos[i] = cSpec.Collaborators[i].GitRepo
	}
	var currContractContent string
	for i := 0; i < noOfCollaborators; i++ {
		if i == 0 {
			currContractContent = GetContractFromRepo(gitRepos[i])
			continue
		}
		intermidiateContractContent := GetContractFromRepo(gitRepos[i])
		if currContractContent != intermidiateContractContent {
			err = fmt.Errorf("contents of Contracts from GitRepo don't match")
			return err
		}
	}
	fmt.Println("Contracts verified successfully")
	return nil
}

func Upload() {

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

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileOptions := &github.RepositoryContentFileOptions{
		Message: github.String("Add example file"),
		Content: data,
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

func Delete() {

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

