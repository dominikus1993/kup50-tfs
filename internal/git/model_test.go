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
		"changeType": "add"
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
		a, ok := change.(map[string]interface{})
		assert.True(t, ok)
		assert.NotNil(t, a["changeType"])
		assert.NotNil(t, a["item"])
	}

}
