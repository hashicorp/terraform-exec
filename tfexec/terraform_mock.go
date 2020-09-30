package tfexec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"time"
)

type MockItemDispenser interface {
	NextMockItem() *MockItem
}

type MockItem struct {
	Args          []string      `json:"args"`
	Stdout        string        `json:"stdout"`
	Stderr        string        `json:"stderr"`
	SleepDuration time.Duration `json:"sleep"`
	ExitCode      int           `json:"exit_code"`

	MockError string `json:"error"`
}

func (m *MockItem) MarshalJSON() ([]byte, error) {
	type t MockItem
	return json.Marshal((*t)(m))
}

func (m *MockItem) UnmarshalJSON(b []byte) error {
	type t MockItem
	return json.Unmarshal(b, (*t)(m))
}

type MockQueue struct {
	Q []*MockItem
}

type MockCall MockItem

func (mc *MockCall) MarshalJSON() ([]byte, error) {
	item := (*MockItem)(mc)
	q := MockQueue{
		Q: []*MockItem{item},
	}
	return json.Marshal(q)
}

func (mc *MockCall) NextMockItem() *MockItem {
	return (*MockItem)(mc)
}

func (mc *MockQueue) NextMockItem() *MockItem {
	if len(mc.Q) == 0 {
		return &MockItem{
			MockError: "no more calls expected",
		}
	}

	var mi *MockItem
	mi, mc.Q = mc.Q[0], mc.Q[1:]

	return mi
}

func NewMockTerraform(md MockItemDispenser) (*Terraform, error) {
	if md == nil {
		md = &MockCall{
			MockError: "no mocks provided",
		}
	}

	wd := os.TempDir()
	tf, err := NewTerraform(wd, os.Args[0])
	if err != nil {
		return nil, err
	}
	tf.mockData = md

	return tf, nil
}

func mockCommand(cmd *exec.Cmd, md MockItemDispenser) *exec.Cmd {
	mockData, err := md.NextMockItem().MarshalJSON()
	if err != nil {
		panic(err)
	}
	
	cmd.Env = append(cmd.Env, "TF_LS_MOCK="+ string(mockData))
	return cmd
}

func ExecuteMockData(rawMockData string) int {
	mi := &MockItem{}
	err := mi.UnmarshalJSON([]byte(rawMockData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to unmarshal mock response: %s", err)
		return 1
	}
	return validateMockItem(mi, os.Args[1:], os.Stdout, os.Stderr)
}

func validateMockItem(m *MockItem, args []string, stdout, stderr io.Writer) int {
	if m.MockError != "" {
		fmt.Fprintf(stderr, m.MockError)
		return 1
	}

	givenArgs := args
	if !reflect.DeepEqual(m.Args, givenArgs) {
		fmt.Fprintf(stderr,
			"arguments don't match.\nexpected: %q\ngiven: %q\n",
			m.Args, givenArgs)
		return 1
	}

	if m.SleepDuration > 0 {
		time.Sleep(m.SleepDuration)
	}

	fmt.Fprint(stdout, m.Stdout)
	fmt.Fprint(stderr, m.Stderr)

	return m.ExitCode
}
