# tablecli

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

Made with ﵆ 
