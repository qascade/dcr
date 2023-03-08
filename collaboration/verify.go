package collaboration

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/go-github/v35/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func GetContractFromRepo(link string) (*string, *error) {

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
		return nil, &err
	}

	content, err := fileContents.GetContent()
	if err != nil {
		err = errors.New("error getting file content from Contract")
		return nil, &err
	}

	log.Println(content)
	return &content, nil
}

func Verify(path string) error {
	err := godotenv.Load("../.env")
	if err != nil {
		err = errors.New("error loading environment variables file")
		return err
	}

	var collabPkg CollaborationParser = &CollaborationPackage{}
	cSpec, _, err := collabPkg.Parse(path)
	if err != nil {
		return errors.New("error parsing contract.yaml file")
	}

	var gitRepos [100]string
	var noOfCollaborators = len(cSpec.Collaborators)

	for i := 0; i < noOfCollaborators; i++ {
		gitRepos[i] = cSpec.Collaborators[i].GitRepo
	}

	var currContractContent *string
	for i := 0; i < noOfCollaborators; i++ {
		if i == 0 {
			currContractContent, _ = GetContractFromRepo(gitRepos[i])
			continue
		}
		intermidiateContractContent, _ := GetContractFromRepo(gitRepos[i])
		if currContractContent != intermidiateContractContent {
			err = errors.New("contents of Contracts from GitRepo don't match")
			return err
		}
	}
	log.Println("Contracts verified successfully")

	return nil
}

func Upload(linkToContractFile string, RepoName string) error {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading environment variables file")
		return errors.New("Error loading environment variables file")
	}

	ctx := context.Background()
	token := os.Getenv("token")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoName := RepoName
	repoDescription := "This is a new repo created by the client"
	repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		Name:        &repoName,
		Description: &repoDescription,
	})
	if err != nil {
		log.Printf("Failed to create repo: %v\n", err)
		return errors.New("Failed to create repo.")
	}

	log.Printf("Created new repo: %v\n", *repo.HTMLURL)

	owner := os.Getenv("owner")
	repoT := RepoName
	filePath := linkToContractFile

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	fileOptions := &github.RepositoryContentFileOptions{
		Message: github.String("Add the contract file of the client"),
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
		log.Println(err)
		return err
	}

	log.Println("File uploaded successfully")

	return nil
}

func Delete(repoName string) error {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading environment variables file")
		return errors.New("Error loading environment variables file")
	}

	owner := os.Getenv("owner")
	repo := repoName
	token := os.Getenv("token")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, er := client.Repositories.Delete(ctx, owner, repo)
	if er != nil {
		log.Println("Error deleting repository:", er)
		return er
	}

	log.Println("Repository deleted successfully")
	return nil
}
