package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const RepoURL = "https://api.github.com/search/repositories"

type RepoSearchResult struct {
	TotalCount       int     `json:"total_count"`
	IncompleteResult bool    `json:"incomplete_result"`
	Repo             []*Repo `json:"items"`
}

type Repo struct {
	Id              int
	NodeId          string `json:"node_id"`
	Name            string
	FullName        string `json:"full_name"`
	Private         bool
	HtmlURL         string `json:"html_url"`
	Description     string
	Fork            bool
	URL             string    `json:"url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Homepage        string    `json:"homepage,emitempty"`
	Size            int
	StargazersCount int `json:"stargazers_count"`
	WatchersCount   int `json:"watchers_count"`
	Language        string
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	MasterBranch    string `json:"master_branch"`
	DefaultBranch   string `json:"default_branch"`
	Score           float32
	Owner           *Owner
}
type Owner struct {
	Login             string
	Id                int
	NodeId            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	URL               string `json:"url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string
}

type byTime []*Repo

func (t byTime) Len() int           { return len(t) }
func (t byTime) Less(i, j int) bool { return t[i].UpdatedAt.Sub(t[j].UpdatedAt) < 0 }
func (t byTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func main() {
	fmt.Println("Please enter the keywords for searching Repositories: ")
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	resp, _ := SearchRepo(strings.Split(string(line), " "))
	var LH []*Repo
	var H []*Repo
	var MH []*Repo
	for _, repo := range (*resp).Repo {
		h := time.Now().Sub(repo.UpdatedAt).Hours()
		if h < 1 {
			LH = append(LH, repo)
		} else if h > 1 && h < 2 {
			H = append(H, repo)
		} else {
			MH = append(MH, repo)
		}
	}
	sort.Sort(byTime(LH))
	sort.Sort(byTime(H))
	sort.Sort(byTime(MH))
	fmt.Println("Update More than Two ago")
	for i, v := range MH {
		fmt.Printf("\t#%2d %-10d %-30s %v\n", i+1, v.Id, v.Name, v.UpdatedAt.Local())
		fmt.Printf("\t\tURL: %s\n", v.URL)
	}
	fmt.Println("Update an hour ago")
	for i, v := range H {
		fmt.Printf("\t#%2d %-10d %-30s %v\n", i+1, v.Id, v.Name, v.UpdatedAt.Local())
		fmt.Printf("\t\tURL: %s\n", v.URL)
	}
	fmt.Println("Update Less than an hour ago")
	for i, v := range LH {
		fmt.Printf("\t#%2d %-10d %-30s %v\n", i+1, v.Id, v.Name, v.UpdatedAt.Local())
		fmt.Printf("\t\tURL: %s\n", v.URL)
	}
}

//SearchRepository queries the Github Repo tracker
func SearchRepo(terms []string) (*RepoSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	if q != "" {
		q = "?q=topic:" + q
	}
	resp, err := http.Get(RepoURL + q)
	if err != nil {
		return nil, err
	}
	//check the status code
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query Faild %s", resp.Status)
	}
	var result RepoSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
