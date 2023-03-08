package git

import (
	"context"
	"errors"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
)

type AzureDevopsClient struct {
	organizationUrl string
	token           string
	project         string
	gitClient       git.Client
}

func NewAzureDevopsClient(ctx context.Context, organizationUrl, token, project string) (*AzureDevopsClient, error) {
	connection := azuredevops.NewPatConnection(organizationUrl, token)
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	return &AzureDevopsClient{organizationUrl: organizationUrl, token: token, gitClient: gitClient}, nil
}

func (client *AzureDevopsClient) GetCommits(ctx context.Context, author string) error {
	repositories, err := client.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{Project: &client.project})
	if err != nil {
		return err
	}
	for _, repo := range *repositories {
		repoId := repo.Id.String()
		commits, commitErr := client.gitClient.GetCommits(ctx, git.GetCommitsArgs{RepositoryId: &repoId, Project: &client.project, SearchCriteria: &git.GitQueryCommitsCriteria{Author: author}})
		if err != nil {
			err = errors.Join(err, commitErr)
			continue
		}
		for _, commit := range *commits {

		}
	}
}
