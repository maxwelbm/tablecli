package tablecli

import (
  "bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	DefaultPadding                        = 2
	DefaultWriter               io.Writer = os.Stdout
	DefaultHeaderFormatter      Formatter
	DefaultFirstColumnFormatter Formatter
	DefaultWidthFunc            WidthFunc = utf8.RuneCountInString
	WidthPersist                []int
)

type Formatter func(string, ...interface{}) string

type WidthFunc func(string) int

type Table interface {
	WithHeaderFormatter(f Formatter) Table
	WithFirstColumnFormatter(f Formatter) Table
	WithPadding(p int) Table
	WithWriter(w io.Writer) Table
	WithWidthFunc(f WidthFunc) Table
	AddRow(vals ...interface{}) Table
	SetRows(Rows [][]string) Table
	Print()
	GetByteFormat() []byte
	CalculateWidths([]string)
	GetHeader() []string
	GetRows() [][]string
	GetWriter() io.Writer 
	PrintHeader(format string)
	PrintRow(format string, row []string)
}

func New(columnHeaders ...interface{}) Table {
	t := table{
		Header: make([]string, len(columnHeaders)),
	}

	t.WithPadding(DefaultPadding)
	t.WithWriter(DefaultWriter)
	t.WithHeaderFormatter(DefaultHeaderFormatter)
	t.WithFirstColumnFormatter(DefaultFirstColumnFormatter)
	t.WithWidthFunc(DefaultWidthFunc)

	for i, col := range columnHeaders {
		t.Header[i] = fmt.Sprint(col)
	}

	return &t
}

type table struct {
	FirstColumnFormatter Formatter
	HeaderFormatter      Formatter
	Padding              int
	Writer               io.Writer
	Width                WidthFunc

	Header []string
	Rows   [][]string
	Widths []int
}

func (t *table) GetRows() [][]string {
	return t.Rows
}

func (t *table) GetWriter() io.Writer {
	return t.Writer
}

func (t *table) GetHeader() []string {
	return t.Header
}

func (t *table) WithHeaderFormatter(f Formatter) Table {
	t.HeaderFormatter = f
	return t
}

func (t *table) WithFirstColumnFormatter(f Formatter) Table {
	t.FirstColumnFormatter = f
	return t
}

func (t *table) WithPadding(p int) Table {
	if p < 0 {
		p = 0
	}

	t.Padding = p
	return t
}

func (t *table) WithWriter(w io.Writer) Table {
	if w == nil {
		w = os.Stdout
	}

	t.Writer = w
	return t
}

func (t *table) WithWidthFunc(f WidthFunc) Table {
	t.Width = f
	return t
}

func (t *table) AddRow(vals ...interface{}) Table {
	maxNumNewlines := 0
	for _, val := range vals {
		maxNumNewlines = max(strings.Count(fmt.Sprint(val), "\n"), maxNumNewlines)
	}
	for i := 0; i <= maxNumNewlines; i++ {
		row := make([]string, len(t.Header))
		for j, val := range vals {
			if j >= len(t.Header) {
				break
			}
			v := strings.Split(fmt.Sprint(val), "\n")
			row[j] = safeOffset(v, i)
		}
		t.Rows = append(t.Rows, row)
	}

	return t
}

func (t *table) SetRows(Rows [][]string) Table {
	t.Rows = [][]string{}
	headerLength := len(t.Header)

	for _, row := range Rows {
		if len(row) > headerLength {
			t.Rows = append(t.Rows, row[:headerLength])
		} else {
			t.Rows = append(t.Rows, row)
		}
	}

	return t
}

func (t *table) Print() {
	format := strings.Repeat("%s", len(t.Header)) + "\n"
	t.CalculateWidths([]string{})

	t.PrintHeader(format)
	for _, row := range t.Rows {
		t.PrintRow(format, row)
	}
}

func (t *table) GetByteFormat() []byte {
  var b bytes.Buffer
	format := strings.Repeat("%s", len(t.Header)) + "\n"
	t.CalculateWidths([]string{})
	for _, row := range t.Rows {
	  vals := t.applyWidths(row, t.Widths)
  	if t.FirstColumnFormatter != nil {
		  vals[0] = t.FirstColumnFormatter("%s", vals[0])
	  }
    b.Write([]byte(fmt.Sprintf(format, vals...)))
	}
  return b.Bytes()
}

func (t *table) PrintHeader(format string) {
	vals := t.applyWidths(t.Header, t.Widths)
	if t.HeaderFormatter != nil {
		txt := t.HeaderFormatter(format, vals...)
		fmt.Fprint(t.Writer, txt)
	} else {
		fmt.Fprintf(t.Writer, format, vals...)
	}
}

func (t *table) PrintRow(format string, row []string) {
	vals := t.applyWidths(row, t.Widths)
	if t.FirstColumnFormatter != nil {
		vals[0] = t.FirstColumnFormatter("%s", vals[0])
	}
	fmt.Fprintf(t.Writer, format, vals...)
}

func (t *table) CalculateWidths(h []string) {
	if len(h) == 0 {
		h = t.Header
	}

	t.Widths = make([]int, len(h))
	for _, row := range t.Rows {
		for i, v := range row {
			if w := t.Width(v) + t.Padding; w > t.Widths[i] {
				t.Widths[i] = w
			}
		}
	}

	for i, v := range t.Header {
		if w := t.Width(v) + t.Padding; w > t.Widths[i] {
			t.Widths[i] = w
		}
	}

	if len(WidthPersist) > 0 {
		for i := 0; i < len(t.Widths); i++ {
			if t.Widths[i] < WidthPersist[i] {
				t.Widths[i] = WidthPersist[i]
			}
		}
	} else {
		WidthPersist = t.Widths
	}
}

func (t *table) applyWidths(row []string, Widths []int) []interface{} {
	out := make([]interface{}, len(row))
	for i, s := range row {
		out[i] = s + t.lenOffset(s, Widths[i])
	}
	return out
}

func (t *table) lenOffset(s string, w int) string {
	l := w - t.Width(s)
	if l <= 0 {
		return ""
	}
	return strings.Repeat(" ", l)
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

func safeOffset(sarr []string, idx int) string {
	if idx >= len(sarr) {
		return ""
	}
	return sarr[idx]
}
