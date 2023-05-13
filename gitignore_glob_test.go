package FastGlobbing_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mxmauro/FastGlobbing"
)

// -----------------------------------------------------------------------------

type GitWildcardMatchTest struct {
	Pattern     string `json:"pattern"`
	Path        string `json:"path"`
	ShouldMatch bool   `json:"should_match"`
}

func TestGitWildcardMatch(t *testing.T) {
	data, err := os.ReadFile("./testdata/gitwildcardmatch_test.json")
	if err != nil {
		t.Fatalf("Error reading test file. Details: %v.", err.Error())
	}

	var tests []GitWildcardMatchTest
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON. Details: %v.", err.Error())
	}

	if runtime.GOOS == "windows" {
		// On Windows, convert paths
		for idx := range tests {
			s := filepath.FromSlash(tests[idx].Path)
			if strings.HasPrefix(s, "\\") {
				s = "C:" + s
			}
			tests[idx].Path = filepath.FromSlash(s)

			s = filepath.FromSlash(tests[idx].Pattern)
			if strings.HasPrefix(s, "\\") {
				s = "C:" + s
			}
			tests[idx].Pattern = filepath.FromSlash(s)
		}
	}

	// Print the list of tests
	for _, test := range tests {
		glob := FastGlobbing.NewGitWildcardMatcher(test.Pattern)

		if glob.Test(test.Path) != test.ShouldMatch {
			var s string

			if test.ShouldMatch {
				s = "match"
			} else {
				s = "mismatch"
			}
			t.Fatalf("Failed to %s pattern '%s' with path '%s'.", s, test.Pattern, test.Path)
		}
	}
}
