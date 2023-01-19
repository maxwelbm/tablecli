package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockClient struct {
	StatusCode int
	Body       []byte
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewReader(m.Body)),
	}, nil
}

func TestOutput(t *testing.T) {
	t.Run("list gists", func(t *testing.T) {
		j, _ := os.ReadFile("./fixtures/gists.json")

		f, _, _ := NewFactory(&MockClient{
			StatusCode: 200,
			Body:       j,
		})
		cmd := List(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
