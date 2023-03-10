package git

import (
	"encoding/json"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/stretchr/testify/assert"
)

var jsonS string = `
{
	"changeCounts": {
	  "Add": 456
	},
	"changes": [
	  {
		"item": {
		  "gitObjectType": "blob",
		  "path": "/MyWebSite/MyWebSite/favicon.ico",
		  "url": "https://dev.azure.com/fabrikam/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/items/MyWebSite/MyWebSite/favicon.ico?versionType=Commit"
		},
		"changeType": "add"
	  },
	  {
		"item": {
		  "gitObjectType": "tree",
		  "path": "/MyWebSite/MyWebSite/fonts",
		  "isFolder": true,
		  "url": "https://dev.azure.com/fabrikam/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/items/MyWebSite/MyWebSite/fonts?versionType=Commit"
		},
		"changeType": "edit"
	  }
	]
  }
`

func TestJsonParsing(t *testing.T) {
	var resp git.GitCommitChanges
	if err := json.Unmarshal([]byte(jsonS), &resp); err != nil {
		panic(err)
	}
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Changes)
	changes := *resp.Changes
	for _, change := range changes {
		subject, err := parseJson(change)
		assert.NoError(t, err)
		assert.NotNil(t, subject)
	}
}

func TestFromJson(t *testing.T) {
	var resp git.GitCommitChanges
	if err := json.Unmarshal([]byte(jsonS), &resp); err != nil {
		panic(err)
	}
	assert.NotNil(t, resp)
	subject, err := FromJson(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, subject)
	assert.Len(t, subject, 2)
}

func TestFilterBlob(t *testing.T) {
	var resp git.GitCommitChanges
	if err := json.Unmarshal([]byte(jsonS), &resp); err != nil {
		panic(err)
	}
	assert.NotNil(t, resp)
	changes, err := FromJson(&resp)

	subject := FilterBlob(changes)
	assert.NoError(t, err)
	assert.NotEmpty(t, subject)
	assert.Len(t, subject, 1)
}

func TestFilterChangeType(t *testing.T) {
	var resp git.GitCommitChanges
	if err := json.Unmarshal([]byte(jsonS), &resp); err != nil {
		panic(err)
	}
	assert.NotNil(t, resp)
	changes, err := FromJson(&resp)

	subject := FilterChangeType(changes, "add")
	assert.NoError(t, err)
	assert.NotEmpty(t, subject)
	assert.Len(t, subject, 1)
}
