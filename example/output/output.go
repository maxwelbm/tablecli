package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	table "github.com/MaxwelMazur/tablecli"
	"github.com/MaxwelMazur/tablecli/example/output/httpmock"
	"github.com/spf13/cobra"
)

type StrList struct {
	ID   string
	Name string
}

type Factory struct {
	HttpClient *http.Client
	IOStreams  *IOStreams
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

func NewFactory(mock *httpmock.Registry) (factory *Factory, out *bytes.Buffer, err *bytes.Buffer) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	f := &Factory{
		HttpClient: &http.Client{Transport: mock},
		IOStreams: &IOStreams{
			Out: stdout,
			Err: stderr,
		},
	}
	return f, stdout, stderr
}

func main() {
	Execute()
}

func Execute() {
	streams := System()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // TODO: Configure this somewhere
	}

	factory := &Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
	}

	cmd := NewRootCmd(factory)
	cobra.CheckErr(cmd.Execute())
}

func NewRootCmd(f *Factory) *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(List(f))
	return rootCmd
}

func List(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			table.DefaultWriter = f.IOStreams.Out
			tbl := table.New("ID", "Name")

			list := []StrList{
				{"12312", "asdfsdfkj"},
			}

			for _, i := range list {
				tbl.AddRow(i.ID, i.Name)
			}

			tbl.Print()
			return nil
		},
	}

	return cmd
}
