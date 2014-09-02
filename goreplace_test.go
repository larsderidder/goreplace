package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const toFind, toReplace, initialContents = "hey", "hello", "hey world\n"

func TestGoReplaceFile(t *testing.T) {
	tempfile, _ := ioutil.TempFile("", "testgoreplace")
	tempfile.WriteString(initialContents)
	patterns := []string{tempfile.Name()}
	goReplace(toFind, toReplace, patterns)

	replaced, filecontents := fileReplaced(tempfile.Name(), toFind, toReplace)
	if !replaced {
		t.Errorf("Expected %s, not %s, in contents: %s", toReplace, toFind, filecontents)
	}
}

func fileReplaced(filename string, toFind string, toReplace string) (bool, string) {
	file, _ := ioutil.ReadFile(filename)
	filecontents := string(file)
	if strings.Contains(filecontents, toFind) || !strings.Contains(filecontents, toReplace) {
		return false, filecontents
	}
	return true, filecontents
}

func TestGoReplacePattern(t *testing.T) {
	pattern_base := "testgoreplace"
	tempfile, _ := ioutil.TempFile(".", pattern_base)
	defer os.Remove(tempfile.Name())

	tempfile.WriteString(initialContents)
	patterns := []string{pattern_base + "*"}
	goReplace(toFind, toReplace, patterns)

	replaced, filecontents := fileReplaced(tempfile.Name(), toFind, toReplace)
	if !replaced {
		t.Errorf("Expected %s, not %s, in contents: %s", toReplace, toFind, filecontents)
	}
}
