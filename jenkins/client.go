package jenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client contains the information for connecting to a jenkins instance
type Client struct {
	Baseurl  string `json:"base_url"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// ViewsResponse describes a response for Views.
type ViewsResponse struct {
	Views []View `json:"views,omitempty"`
}

// View describes what a view looks like
type View struct {
	Name string `json:"name,omitempty"`
	Jobs []Job  `json:"jobs,omitempty"`
}

// Job describes a job object from the Jenkins API.
type Job struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"fullDisplayName,omitempty"`
	LastBuild   Build  `json:"lastBuild,omitempty"`
}

// Build describes a build from the Jenkins API.
type Build struct {
	Result    string `json:"result,omitempty"`
	Number    int    `json:"number,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// New provides a jenkins client
func New(uri, username, token string) *Client {
	return &Client{
		Baseurl:  uri,
		Username: username,
		Token:    token,
	}
}

// GetJobs gets the jobs for a given Jenkins view.
func (c *Client) GetJobs(jenkinsView string) ([]Job, error) {
	if jenkinsView == "" {
		jenkinsView = "all"
	}

	urlEndpoint := fmt.Sprintf("api/json?tree=%s&depth=1", url.QueryEscape("views[name,jobs[name,fullDisplayName,lastBuild[number,timestamp,result]]]"))
	url := fmt.Sprintf("%s/%s", c.Baseurl, urlEndpoint)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request to get jenkins views %s responded with status %d", urlEndpoint, resp.StatusCode)
	}

	var v ViewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("decoding json response for jenkins views from %s failed: %v", urlEndpoint, err)
	}

	for _, view := range v.Views {
		if view.Name == jenkinsView {
			return view.Jobs, nil
		}
	}

	return nil, fmt.Errorf("unable to get named jenkins view: %s", jenkinsView)
}
