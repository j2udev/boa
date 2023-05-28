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
  options [flags] [options]

Options:
  option1, opt1   opt1 description
  option2         opt2 description

Flags:
  -h, --help   help for options
`
	expectedProfilesOutput := `Usage:
  profiles [flags] [options]

Options:
  option1, opt1   opt1 description
  option2         opt2 description

Profiles:
  profile1, prof1   prof1 description
    ↳ Options:      opt1, option2
  profile2          prof2 description
    ↳ Options:      option1

Flags:
  -h, --help   help for profiles
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
	profiles := []Profile{
		{
			Args: []string{"profile1, prof1"},
			Opts: []string{"opt1", "option2"},
			Desc: "prof1 description",
		},
		{
			Args: []string{"profile2"},
			Opts: []string{"option1"},
			Desc: "prof2 description",
		},
	}

	cmd1 := NewCmd("options").
		WithOptions(options...).
		WithOptionsTemplate().
		WithNoOp().
		Build()
	cmd2 := NewCmd("profiles").
		WithOptions(options...).
		WithProfiles(profiles...).
		WithOptionsTemplate().
		WithNoOp().
		Build()

	assert.Equal(t, expectedOptionsOutput, captureCmdOutput(cmd1, "-h"))
	assert.Equal(t, expectedProfilesOutput, captureCmdOutput(cmd2, "-h"))
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
