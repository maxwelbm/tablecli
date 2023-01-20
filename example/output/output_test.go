package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var expected string = "ID                                URL                                                            \n85cbedcda90fd5e4e8fdd9523ade5724  https://api.github.com/gists/85cbedcda90fd5e4e8fdd9523ade5724  \ne23a1e6db7006e89b6dbb90c1ae24bb4  https://api.github.com/gists/e23a1e6db7006e89b6dbb90c1ae24bb4  \nd8d6f864291944fbc9b2de93a5587a49  https://api.github.com/gists/d8d6f864291944fbc9b2de93a5587a49  \na06fe81926065a828c22ea74bd00bea1  https://api.github.com/gists/a06fe81926065a828c22ea74bd00bea1  \nbc972ef81fef5251efe3c17fe417e552  https://api.github.com/gists/bc972ef81fef5251efe3c17fe417e552  \na7a4f8a97a1b797f3a6f9c9bf9817cf0  https://api.github.com/gists/a7a4f8a97a1b797f3a6f9c9bf9817cf0  \n06fac34d4c84ccda17f098f633937cb7  https://api.github.com/gists/06fac34d4c84ccda17f098f633937cb7  \n1a1802ff167b84bb6afcded3304b004c  https://api.github.com/gists/1a1802ff167b84bb6afcded3304b004c  \nefcfe694816756896032ef8b41375bac  https://api.github.com/gists/efcfe694816756896032ef8b41375bac  \n8c42c281085a2f75fc5d985f1a46985f  https://api.github.com/gists/8c42c281085a2f75fc5d985f1a46985f  \nee22500427c6c52cb8d84376110d7339  https://api.github.com/gists/ee22500427c6c52cb8d84376110d7339  \nc236a4da659a6adb3261205433ff6328  https://api.github.com/gists/c236a4da659a6adb3261205433ff6328  \n4da1c2bc6140e12f228f3f362bb9c124  https://api.github.com/gists/4da1c2bc6140e12f228f3f362bb9c124  \nc35351a29f8bdcad0443f4f7a379f8b2  https://api.github.com/gists/c35351a29f8bdcad0443f4f7a379f8b2  \n4129197ee55bfcbd90858f330a3bb746  https://api.github.com/gists/4129197ee55bfcbd90858f330a3bb746  \ndadd6f8bc79a1c05d7513a419d705c8c  https://api.github.com/gists/dadd6f8bc79a1c05d7513a419d705c8c  \n"

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
		t.Parallel()

		j, _ := os.ReadFile("./fixtures/gists.json")

		factory := NewFactory(&MockClient{
			StatusCode: 200,
			Body:       j,
		})

		cmd := List(factory)
		cmd.SetArgs([]string{})
		_, err := cmd.ExecuteC()

		require.NoError(t, err)
		require.Equal(t, expected, factory.Buffer.String())
	})
}
