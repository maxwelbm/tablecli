package main

import (
	"testing"

	"github.com/MaxwelMazur/tablecli/example/output/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	t.Run("list page 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications"),
			httpmock.JSONFromFile("./fixtures/applications.json"),
		)

		f, stdout, _ := NewFactory(mock)
		cmd := List(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID     Name       \n12312  asdfsdfkj  \n", stdout.String())
	})
}
