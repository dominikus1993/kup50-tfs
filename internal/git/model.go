package git

import (
	"errors"

	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/samber/lo"
)

var errCantParse = errors.New("cant parse gitChange from json")
var errCantItemParse = errors.New("cant parse item for gitChange from json")

type gitChange struct {
	Item struct {
		GitObjectType string `json:"gitObjectType"`
		Path          string `json:"path"`
		URL           string `json:"url"`
	} `json:"item"`
	ChangeType string `json:"changeType"`
}

func FromJson(change *git.GitCommitChanges) ([]*gitChange, error) {
	if change.Changes == nil {
		return make([]*gitChange, 0), nil
	}
	result := make([]*gitChange, len(*change.Changes))
	changes := *change.Changes
	var errorRes error
	for i, v := range changes {
		obj, err := parseJson(v)
		if err != nil {
			errorRes = errors.Join(errorRes, err)
		}
		result[i] = obj
	}
	return result, errorRes
}

func FilterChangeType(changes []*gitChange, changeTypes ...string) []*gitChange {
	if len(changeTypes) == 0 {
		return changes
	}
	return lo.Filter(changes, func(change *gitChange, _ int) bool {
		return lo.Contains(changeTypes, change.ChangeType)
	})
}

func FilterBlob(changes []*gitChange) []*gitChange {
	return lo.Filter(changes, func(change *gitChange, _ int) bool {
		return change.Item.GitObjectType == "blob"
	})
}

func parseJson(json interface{}) (*gitChange, error) {
	object, ok := json.(map[string]interface{})
	if !ok {
		return nil, errCantParse
	}
	item, ok := object["item"].(map[string]interface{})
	if !ok {
		return nil, errCantItemParse
	}
	result := gitChange{}
	result.ChangeType = object["changeType"].(string)
	result.Item.GitObjectType = item["gitObjectType"].(string)
	result.Item.Path = item["path"].(string)
	result.Item.URL = item["url"].(string)
	return &result, nil
}
