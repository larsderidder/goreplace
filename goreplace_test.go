package main

import ( 
    "testing"
    "strings"
    "io/ioutil"
)

func TestGoReplace(t *testing.T) {
    to_find, to_replace := "hey", "hello"
    contents := "hey world\n"
    tempfile, _ := ioutil.TempFile(".", "testgoreplace")
    tempfile.WriteString(contents)
    patterns := []string{tempfile.Name()}
    goReplace(to_find, to_replace, patterns)
    file, _ := ioutil.ReadFile(tempfile.Name())
    if filecontents := string(file); strings.Contains(filecontents, to_find) || ! strings.Contains(filecontents, to_replace) {
        t.Errorf("Expected %s, not %s, in contents: %s", to_replace, to_find, filecontents)
    }
}
