package git

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dominikus1993/kup50-tfs/internal/datetime"
	"github.com/dominikus1993/kup50-tfs/internal/diff"
	"github.com/dominikus1993/kup50-tfs/internal/dir"
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
	cfg             *diff.BlobDiffer
}

func NewAzureDevopsClient(ctx context.Context, organizationUrl, token, project string, cfg *diff.BlobDiffer) (*AzureDevopsClient, error) {
	connection := azuredevops.NewPatConnection(organizationUrl, token)
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	return &AzureDevopsClient{organizationUrl: organizationUrl, token: token, gitClient: gitClient, cfg: cfg}, nil
}

func (client *AzureDevopsClient) GetChanges(ctx context.Context, author string) <-chan *RepositoryChanges {
	result := make(chan *RepositoryChanges, 200)
	firstDay, lastDay := datetime.FirstAndLastDayOfTheMonth(time.Now())
	repositories, err := client.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{Project: &client.project})
	if err != nil {
		panic(err)
	}
	go func() {
		log.WithField("firstDay", firstDay).WithField("lastDay", lastDay).WithField("project", client.project).Infoln("Start processing")
		for _, repo := range *repositories {
			repoId := repo.Id.String()
			commits, commitErr := client.gitClient.GetCommits(ctx, git.GetCommitsArgs{RepositoryId: &repoId, Project: &client.project, SearchCriteria: &git.GitQueryCommitsCriteria{Author: &author, FromDate: datetime.FormatToAzureDevopsTime(firstDay), ToDate: datetime.FormatToAzureDevopsTime(lastDay)}})
			if commitErr != nil {
				err = errors.Join(err, commitErr)
				log.WithField("repoName", *repo.Name).WithError(err).Warnln("can't download commits")
				continue
			}
			repository := NewRepositoryChanges(&repo)
			for _, commit := range *commits {
				var authorDate = commit.Author.Date
				if authorDate.Time.Before(firstDay) {
					continue
				}
				var author = *commit.Author.Email
				var repoName = *repo.Name
				log.WithField("date", authorDate).WithField("email", author).WithField("repo", repoName).Infoln("commit")
				changes, err := client.gitClient.GetChanges(ctx, git.GetChangesArgs{CommitId: commit.CommitId, Project: &client.project, RepositoryId: &repoId})
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
				log.WithField("repo", *repo.Name).Infoln("not empty")
				result <- repository
				continue
			}
		}
		close(result)
	}()
	return result
}

func (client *AzureDevopsClient) DowloadAndSaveChanges(ctx context.Context, stream <-chan *RepositoryChanges) {
	for repo := range stream {
		log.WithField("repo", repo.repoName).Infoln("Start Saving")
		for _, change := range repo.changes {
			directory, err := dir.CreateDir(repo.repoName)
			if err != nil {
				log.WithField("repoName", repo.repoName).WithError(err).Error("can't create directory")
				continue
			}
			filename := createFilename(*directory, change.Item.Path, change.ChangeType)
			switch change.ChangeType {
			case "add":
				changes, err := client.gitClient.GetBlobContent(ctx, git.GetBlobContentArgs{RepositoryId: &repo.repoId, Project: &client.project, Download: &download, Sha1: &change.Item.ObjectId})

				if err != nil {
					log.WithField("repoName", repo.repoName).WithError(err).Error("can't dowload commit changes blob")
				}
				err = dir.Save(filename, changes)
				// handle err
				if err != nil {
					log.WithField("repoName", repo.repoName).WithField("filename", filename).WithError(err).Error("can't save commit changes blob")
				}
			case "edit":
				newchanges, err := client.gitClient.GetBlobContent(ctx, git.GetBlobContentArgs{RepositoryId: &repo.repoId, Project: &client.project, Download: &download, Sha1: &change.Item.ObjectId})

				if err != nil {
					log.WithField("repoName", repo.repoName).WithField("filename", filename).WithError(err).Error("can't downloadnew commit changes blob")
				}
				oldchanges, err := client.gitClient.GetBlobContent(ctx, git.GetBlobContentArgs{RepositoryId: &repo.repoId, Project: &client.project, Download: &download, Sha1: &change.Item.OriginalObjectId})

				if err != nil {
					log.WithField("repoName", repo.repoName).WithField("filename", filename).WithError(err).Error("can't download old commit changes blob")
				}

				res, err := client.cfg.DiffBlobs(oldchanges, newchanges)
				if err != nil {
					log.WithError(err).WithField("repoName", repo.repoName).Errorln("Can't diff html")
					continue
				}
				// handle err
				err = dir.SaveString(filename, res)
				// handle err
				if err != nil {
					log.WithField("repoName", repo.repoName).WithField("filename", filename).WithError(err).Error("can't save commit changes blob")
				}
			}

		}
	}
}

func createFilename(dir, path, operation string) string {
	return filepath.Join(dir, filepath.Base(fmt.Sprintf("changes_%s_%s.html", strings.ReplaceAll(path, "/", "_"), operation)))
}
