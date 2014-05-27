package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "os"
    "io/ioutil"
    "strings"
    "regexp"
)

func main() {
    verbose := true
    args := os.Args[1:]
    if verbose {
        fmt.Println(os.Args)
    }
    if len(args) < 3 {
        panic("Arguments dude!")
    }
    to_find := args[0]
    to_replace := args[1]
    patterns := args[2:]
    if verbose {
        fmt.Println(to_find)
        fmt.Println(to_replace)
        fmt.Println(patterns)
    }

    ch := make(chan bool)
    files, _ := ioutil.ReadDir(".")
    found := 0
    for _, file := range files {
        if ! file.IsDir() {
            for _, pattern := range patterns {
                match, _ := regexp.Match(pattern, []byte(file.Name()))
                if match {
                    found += 1
                    go findReplace(file.Name(), to_find, to_replace, ch)
                }
            }
        }
    }
    changed := 0
    for i := 0; i < found; i++ {
        result := <-ch
        if result {
            changed += 1
        }
    }
    fmt.Printf("Done! Changed %d files.\n", changed)
}

func findReplace(filename, to_find, to_replace string, ch chan bool) {
    file, _ := ioutil.ReadFile(filename)
    filecontents := string(file)
    if strings.Contains(filecontents, to_find) {
        filecontents = strings.Replace(filecontents, to_find, to_replace, -1)
        ioutil.WriteFile(filename, []byte(filecontents), 0644)
        ch <- true
    } else {
        ch <- false
    }
}

func catFile(filename string) {
    cmd := exec.Command("cat", filename)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Run()
    fmt.Println(out.String())
}
