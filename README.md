# tablecli

Packages provide a convenient way to generate tabular output of any data, useful primarily for CLI tools.

<img src="https://raw.githubusercontent.com/MaxwelMazur/tablecli/main/0aac4e6a54c170b06e2bd3848d2b735e.gif?auto=compress&cs=tinysrgb&h=750&w=1260" alt="Girl in a jacket" width="100%" height="250px">

#### Example of use:
```go 
package main

import (
    table "github.com/MaxwelMazur/tablecli"
    "github.com/fatih/color"
    "strings"
)


type list struct {
    ID string 
    Name string
}

func main() {
    tbl := table.New("ID", "NAME")
    headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgGreen).SprintfFunc()
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

    var list = []list{
        { "123123", "Jonh"},
        { "123121", "Jeff"},
    }

    for _, i := range list {
	tbl.AddRow(i.ID, i.Name)
    }

    format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
    tbl.CalculateWidths([]string{})

    tbl.PrintHeader(format)
	for _, r := range tbl.GetRows() {
	    tbl.PrintRow(format, r)
    }
}

```

#### Output: 
```sh
ID      NAME  
123123  Jonh  
123121  Jeff 
```

[![GoDoc](https://godoc.org/github.com/MaxwelMazur/tablecli?status.svg)](https://godoc.org/github.com/MaxwelMazur/tablecli)<br>

