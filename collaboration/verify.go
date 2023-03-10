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

type GitRepoContent string

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

func (c *CollaborationPackage) Verify(path string) error {
	err := godotenv.Load("../.env")
	if err != nil {
		err = errors.New("error loading environment variables file")
		return err
	}

	var collabPkg Collaboration
	collabPkg, err = NewCollaborationPkg(path)
	if err != nil {
		return errors.New("error parsing contract.yaml file")
	}
	cSpec := collabPkg.GetContractSpec()
	var gitRepos [100]string
	var noOfCollaborators = len(cSpec.Collaborators)

	for i := 0; i < noOfCollaborators; i++ {
		gitRepos[i] = cSpec.Collaborators[i].GitRepo
	}

	var currContractContent *string
	for i := 0; i < noOfCollaborators; i++ {
		if i == 0 {
			var err *error
			currContractContent, err = GetContractFromRepo(gitRepos[i])
			if err != nil {
				return errors.New("error getting contract.yaml file from GitRepo")
			}
			continue
		}
		intermidiateContractContent, err := GetContractFromRepo(gitRepos[i])
		if err != nil {
			return errors.New("error getting contract.yaml file from GitRepo")
		}
		if currContractContent != intermidiateContractContent {
			err := errors.New("contents of Contracts from GitRepo don't match")
			return err
		}
	}
	log.Println("Contracts verified successfully")

	return nil
}

func (c *CollaborationPackage) UploadToRepo(path string) error {

	// TODO - Path needed here is to the contract.yaml file. Not the directory.
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading environment variables file")
		return errors.New("error loading environment variables file")
	}

	ctx := context.Background()
	token := os.Getenv("token")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoDescription := "This is a new repo created by the client"
	repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		Name:        &c.repoName,
		Description: &repoDescription,
	})
	if err != nil {
		log.Printf("Failed to create repo: %v\n", err)
		return errors.New("failed to create repo")
	}

	log.Printf("Created new repo: %v\n", *repo.HTMLURL)

	owner := os.Getenv("owner")
	repoT := c.repoName

	data, err := os.ReadFile(path)
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
		path,
		fileOptions,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("File uploaded successfully")

	return nil
}

func (c *CollaborationPackage) Terminate() error {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading environment variables file")
		return errors.New("error loading environment variables file")
	}

	owner := os.Getenv("owner")
	token := os.Getenv("token")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, err = client.Repositories.Delete(ctx, owner, c.repoName)
	if err != nil {
		log.Println("Error deleting repository:", err)
		return err
	}

	log.Println("Repository deleted successfully")
	return nil
}
