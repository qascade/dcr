# DCR Pull Request Guidelines

## Code Style

This repository uses `gofmt` to maintain code style and consistency. 
Please run `gofmt -s -w .` and `go vet ./...` before pushing the commits. 

Use `go generate ./...` to generate mocks. 

This Repo will incorporate all the recommendations of [effective go doc](https://go.dev/doc/effective_go) and addition to that will be following [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Rules

There are a few basic ground-rules for contributors:

1. Non-main branches should be used for ongoing work.
2. **Pull requests** must pass all the workflows and should have atleast one approval before merging into **main**.  
3. **Commits** should be as atomic as possible. Don't commit before a small milestone is achieved. Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) to log commit messages. Divide the full issue into subtasks that are independent of each other and make a single commit for them. A commit should not have any errors/incomplete or dead code. 
4. We will be following [semantic versioning](https://semver.org/). This along with point will help in auto-generating changelogs. 
5. PR will be closed if commit conventions not followed. 
5. Everything is to be tracked through Github issues. Raise an issue with appropriate tags and link a PR with it. 

### Releases

Declaring formal releases remains the prerogative of the project maintainer.

