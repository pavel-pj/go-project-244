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

	//сначала проверяется относительынй путь в текущей папке

	data, err := os.ReadFile(wd + "/" + "fixtures/diff_json.txt")
	if err != nil {
		t.Error(err)
	}
	//mt.Println(want)
	want := string(data)
	//fmt.Println(want)

	got, err := GendDiff("files/file1.json", "files/file2.json", "stylish")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(got)

	require.Equal(t, want, got)
}
