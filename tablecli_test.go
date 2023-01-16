package tablecli

import(
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	t.Parallel()

	var formatter Formatter

	fn := func(a string, b ...interface{}) string { return "" }
	f := Formatter(fn)

	assert.IsType(t, formatter, f)
}

func TestTable_New(t *testing.T) {
	t.Parallel()

	buf := bytes.Buffer{}
	New("foo", "bar").WithWriter(&buf).Print()
	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "bar")

	buf.Reset()
	New().WithWriter(&buf).Print()
	out = buf.String()

	assert.Empty(t, strings.TrimSpace(out))
}
