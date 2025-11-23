package code

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiffTest(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cases := []struct {
		name, file1, file2, format, expected string
		hasError                             bool
	}{
		{
			name:     "json stylish",
			file1:    "files/file1.json",
			file2:    "files/file2.json",
			format:   "stylish",
			expected: "fixtures/diff_stylish.txt",
			hasError: false,
		},
		{
			name:     "json plain",
			file1:    "files/file1.json",
			file2:    "files/file2.json",
			format:   "plain",
			expected: "fixtures/diff_plain.txt",
			hasError: false,
		},
		{
			name:     "json json",
			file1:    "files/file1.json",
			file2:    "files/file2.json",
			format:   "json",
			expected: "fixtures/diff_json.txt",
			hasError: false,
		},
		{
			name:     "yaml stylish",
			file1:    "files/file1.yaml",
			file2:    "files/file2.yaml",
			format:   "stylish",
			expected: "fixtures/diff_stylish.txt",
			hasError: false,
		},
		{
			name:     "yaml plain",
			file1:    "files/file1.yaml",
			file2:    "files/file2.yaml",
			format:   "plain",
			expected: "fixtures/diff_plain.txt",
			hasError: false,
		},
		{
			name:     "yaml json",
			file1:    "files/file1.yaml",
			file2:    "files/file2.yaml",
			format:   "json",
			expected: "fixtures/diff_json.txt",
			hasError: false,
		},
		{
			name:     "yaml stylish",
			file1:    "files/file1.yml",
			file2:    "files/file2.yml",
			format:   "stylish",
			expected: "fixtures/diff_stylish.txt",
			hasError: false,
		},
		{
			name:     "yaml plain",
			file1:    "files/file1.yml",
			file2:    "files/file2.yml",
			format:   "plain",
			expected: "fixtures/diff_plain.txt",
			hasError: false,
		},
		{
			name:     "yaml json",
			file1:    "files/file1.yml",
			file2:    "files/file2.yml",
			format:   "json",
			expected: "fixtures/diff_json.txt",
			hasError: false,
		},
	}

	for _, r := range cases {

		t.Run(r.name, func(t *testing.T) {

			data, _ := os.ReadFile(wd + "/" + r.expected)
			want := string(data)

			got, _ := GenDiff(r.file1, r.file2, r.format)
			require.Equal(t, want, got)

		})
	}
}
