package git

type gitChange struct {
	Item struct {
		GitObjectType string `json:"gitObjectType"`
		Path          string `json:"path"`
		URL           string `json:"url"`
	} `json:"item"`
	ChangeType string `json:"changeType"`
}
