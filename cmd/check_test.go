package cmd

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	err := os.Chdir("../testdata/testmod")
	if err != nil {
		t.Fatal(err)
	}
	_, err = checkDeps(false)
	if err != nil {
		t.Fatal(err)
	}
}
