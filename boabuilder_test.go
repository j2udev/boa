package boa

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBoaCmdBuilder(t *testing.T) {
	expectedOptionsOutput := `Usage:
  test [flags] [options]

Options:
  option1, opt1   opt1 description
  option2         opt2 description

Flags:
  -h, --help   help for test
`

	options := []Option{
		{
			Args: []string{"option1, opt1"},
			Desc: "opt1 description",
		},
		{
			Args: []string{"option2"},
			Desc: "opt2 description",
		},
	}

	cmd1 := NewCmd("test").
		WithOptions(options...).
		WithOptionsTemplate().
		WithNoOp().
		Build()
	cmd2 := NewCmd("test").
		WithOptionsAndTemplate(options...).
		WithNoOp().
		Build()

	assert.Equal(t, expectedOptionsOutput, captureCmdOutput(cmd1, "-h"))
	assert.Equal(t, expectedOptionsOutput, captureCmdOutput(cmd2, "-h"))
}

func captureCmdOutput(cmd *cobra.Command, args ...string) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmd.SetArgs(args)
	cmd.Execute()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := string(out)
	os.Stdout = rescueStdout
	return output
}
