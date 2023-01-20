package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	table "github.com/MaxwelMazur/tablecli"
	"github.com/spf13/cobra"
)

// HttpClient interface made to fake a request
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Gists []struct {
	URL         string      `json:"url,omitempty"`
	ForksURL    string      `json:"forks_url,omitempty"`
	CommitsURL  string      `json:"commits_url,omitempty"`
	ID          string      `json:"id,omitempty"`
	NodeID      string      `json:"node_id,omitempty"`
	GitPullURL  string      `json:"git_pull_url,omitempty"`
	GitPushURL  string      `json:"git_push_url,omitempty"`
	HTMLURL     string      `json:"html_url,omitempty"`
	Public      bool        `json:"public,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	Description string      `json:"description,omitempty"`
	Comments    int         `json:"comments,omitempty"`
	User        interface{} `json:"user,omitempty"`
	CommentsURL string      `json:"comments_url,omitempty"`
	Truncated   bool        `json:"truncated,omitempty"`
}

type Factory struct {
	HttpClient HttpClient
	IOStreams  *IOStreams
	Buffer     bytes.Buffer
}

type IOStreams struct {
	In  io.ReadCloser
	Out io.Writer
	Err io.Writer
}

func System() *IOStreams {
	return &IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

func NewFactory(client HttpClient) *Factory {
	return &Factory{
		HttpClient: client,
		IOStreams:  System(),
		Buffer:     bytes.Buffer{},
	}
}

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // TODO: Configure this somewhere
	}

	factory := NewFactory(httpClient)
	rootCmd := new(cobra.Command)
	rootCmd.AddCommand(List(factory))
	cobra.CheckErr(rootCmd.Execute())
}

func List(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {

			tbl := table.New("ID", "URL").WithWriter(&f.Buffer)
			for _, i := range getGists(f) {
				tbl.AddRow(i.ID, i.URL)
			}

			tbl.Print()

			return nil
		},
	}
	return cmd
}

//	curl \
//	  -H "Accept: application/vnd.github+json" \
//	  -H "Authorization: Bearer <YOUR-TOKEN>"\
//	  -H "X-GitHub-Api-Version: 2022-11-28" \
//	  https://api.github.com/users/<YOUR-USER>/gists
func getGists(f *Factory) Gists {
	requestURL := "https://api.github.com/users/<YOUR-USER>/gists"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer <YOUR-TOKEN>") // replace <YOUR-TOKEN> for you token
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := f.HttpClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var gists Gists

	if err := json.Unmarshal(resBody, &gists); err != nil {
		fmt.Printf("client: error Unmarshal: %s\n", err)
		os.Exit(1)
	}
	return gists
}
