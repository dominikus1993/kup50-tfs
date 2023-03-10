package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dominikus1993/kup50-tfs/internal/datetime"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	log "github.com/sirupsen/logrus"
)

var download = true

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

func (client *AzureDevopsClient) GetChanges(ctx context.Context, author string) <-chan *RepositoryChanges {
	result := make(chan *RepositoryChanges)
	firstDay, lastDay := datetime.FirstAndLastDayOfTheMonth(time.Now())
	repositories, err := client.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{Project: &client.project})
	if err != nil {
		panic(err)
	}
	go func() {
		for _, repo := range *repositories {
			repoId := repo.Id.String()
			commits, commitErr := client.gitClient.GetCommits(ctx, git.GetCommitsArgs{RepositoryId: &repoId, Project: &client.project, SearchCriteria: &git.GitQueryCommitsCriteria{Author: &author, FromDate: &firstDay, ToDate: &lastDay}})
			if commitErr != nil {
				err = errors.Join(err, commitErr)
				log.WithField("repoName", *repo.Name).WithError(err).Warnln("can't download commits")
				continue
			}
			repository := NewRepositoryChanges(&repo)
			for _, commit := range *commits {
				changes, err := client.gitClient.GetChanges(ctx, git.GetChangesArgs{CommitId: commit.CommitId, RepositoryId: &repoId})
				if err != nil {
					err = errors.Join(err, commitErr)
					log.WithError(err).Warnln("can't download changes")
					continue
				}
				gitchanges, err := FromJson(changes)
				if err != nil {
					err = errors.Join(err, commitErr)
					log.WithError(err).Warnln("can't parse changes")
					continue
				}

				filteredChanges := FilterChangeType(FilterBlob(gitchanges), "add", "edit")

				repository.AddChanges(filteredChanges)
			}
			if repository.HasChanges() {
				result <- repository
			}
		}
		close(result)
	}()
	return result
}

func (client *AzureDevopsClient) DowloadAndSaveChanges(ctx context.Context, stream <-chan *RepositoryChanges) {
	for repo := range stream {
		for _, change := range repo.changes {
			switch change.ChangeType {
			case "add":
				changes, err := client.gitClient.GetBlobContent(ctx, git.GetBlobContentArgs{RepositoryId: repo.repoId, Project: &client.project, Download: &download, Sha1: &change.Item.ObjectId})

				if err != nil {
					log.WithField("repoName", *repo.repo.Name).WithError(err).Error("can't dowload commit changes blob")
				}
				outFile, _ := os.Create(fmt.Sprintf("changes_%s_%s.html", *repo.repoId, change.Item.ObjectId))
				// handle err
				_, err = io.Copy(outFile, changes)
				// handle err
				if err != nil {
					log.WithField("repoName", *repo.repo.Name).WithError(err).Error("can't dowload commit changes blob")
				}
				outFile.Close()
			case "edit":
				changes, err := client.gitClient.GetBlobContent(ctx, git.GetBlobContentArgs{RepositoryId: repo.repoId, Project: &client.project, Download: &download, Sha1: &change.Item.ObjectId})

				if err != nil {
					log.WithField("repoName", *repo.repo.Name).WithError(err).Error("can't dowload commit changes blob")
				}
				outFile, _ := os.Create(fmt.Sprintf("jp2137_%s_%s.html", *repo.repoId, change.Item.ObjectId))
				// handle err
				_, err = io.Copy(outFile, changes)
				// handle err
				if err != nil {
					log.WithField("repoName", *repo.repo.Name).WithError(err).Error("can't dowload commit changes blob")
				}
				outFile.Close()
			}
		}
	}
}
