package tfexec

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseWorkspaceList(t *testing.T) {
	for i, c := range []struct {
		expected        []string
		expectedCurrent string
		stdout          string
	}{
		{
			[]string{"default"},
			"default",
			`* default

`,
		},
		{
			[]string{"default", "foo", "bar"},
			"foo",
			`  default
* foo
  bar

`,
		},

		// linux new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\n* foo\n\n",
		},
		// windows new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\r\n* foo\r\n\r\n",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actualList, actualCurrent := parseWorkspaceList(c.stdout)

			if actualCurrent != c.expectedCurrent {
				t.Fatalf("expected selected %q, got %q", c.expectedCurrent, actualCurrent)
			}

			if !reflect.DeepEqual(c.expected, actualList) {
				t.Fatalf("expected %#v, got %#v", c.expected, actualList)
			}
		})
	}
}
